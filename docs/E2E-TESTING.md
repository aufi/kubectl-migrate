# E2E Testing Guide for kubectl-migrate

Tento dokument popisuje E2E testovací framework pro kubectl-migrate CLI nástroj.

## Přehled

E2E testovací framework poskytuje kompletní infrastrukturu pro testování CLI příkazů proti reálnému Kubernetes clusteru.

### Klíčové vlastnosti

- ✅ **CLI Execution** - Spouštění a validace CLI příkazů
- ✅ **Cluster Management** - Interakce s Kubernetes clusterem
- ✅ **Assertions** - Bohaté assertion helpers pro validaci
- ✅ **Test Isolation** - Každý test běží v vlastním namespace s automatickým cleanupem
- ✅ **Sample Resources** - Využití existujících testovacích aplikací
- ✅ **CI/CD Integration** - GitHub Actions workflow

## Struktura

```
test/e2e/
├── framework/              # Testovací framework
│   ├── cli.go             # CLI executor pro spouštění příkazů
│   ├── cluster.go         # Kubernetes cluster management
│   ├── assertions.go      # Assertion helpers
│   └── suite.go           # Test suite pomocníci
├── export_test.go         # E2E testy pro export příkaz
├── README.md              # Uživatelská dokumentace
└── CONTRIBUTING.md        # Příspěvatelská dokumentace
```

## Rychlý start

### 1. Příprava prostředí

```bash
# Vytvořte Kubernetes cluster (např. Kind)
kind create cluster --name e2e-test

# Nebo použijte existující cluster
kubectl cluster-info
```

### 2. Spuštění testů

```bash
# Build binárky
make build

# Spusťte všechny E2E testy
make test-e2e

# Spusťte pouze export testy
make test-e2e-export

# Rychlé testy (přeskočí dlouhé testy)
make test-e2e-quick
```

### 3. Spuštění konkrétního testu

```bash
# Spusťte konkrétní test
E2E_BINARY=./bin/kubectl-migrate go test -v \
  ./test/e2e/export_test.go \
  ./test/e2e/framework/*.go \
  -run TestExportCommand/basic
```

## Makefile příkazy

| Příkaz | Popis |
|--------|-------|
| `make test-e2e` | Spustí všechny E2E testy (vyžaduje cluster) |
| `make test-e2e-quick` | Spustí rychlé E2E testy (short mode) |
| `make test-e2e-export` | Spustí pouze export E2E testy |
| `make test-unit` | Spustí unit testy |
| `make test-all` | Spustí unit + E2E testy |

## Framework komponenty

### CLI Executor

Spouští CLI příkazy a zachytává výstup:

```go
// Základní spuštění
result := suite.RunCLI("export", "--namespace", "default")

// Kontrola výsledku
if result.Success() {
    fmt.Println("Příkaz uspěl")
}

// Přístup k výstupu
fmt.Println(result.Stdout)
fmt.Println(result.Stderr)
fmt.Println(result.ExitCode)
```

### Cluster Manager

Interakce s Kubernetes clusterem:

```go
// Čekání na deployment
err := suite.Cluster.WaitForDeployment("default", "my-app", 120*time.Second)

// Získání podů
pods, err := suite.Cluster.GetPods("default", map[string]string{"app": "nginx"})

// Kontrola existence resource
exists, err := suite.Cluster.ResourceExists(gvr, "default", "my-resource")
```

### Assertions

Validace výsledků testů:

```go
// Příkazové assertions
suite.Assert.AssertCommandSuccess(result)
suite.Assert.AssertOutputContains(result, "expected text")

// Souborové assertions
suite.Assert.AssertFileExists("/path/to/file")
suite.Assert.AssertYAMLFileValid("/path/to/file.yaml")

// Cluster assertions
suite.Assert.AssertResourceExists(cluster, gvr, "namespace", "name")
```

## Psaní testů

### Základní struktura testu

```go
func TestMyFeature(t *testing.T) {
    suite := framework.NewTestSuite(t)
    defer suite.Cleanup()

    // Přeskočit v short mode
    suite.SkipIfShort("Vyžaduje cluster")

    // Vytvořit namespace
    ns := suite.CreateTestNamespace("my-test")

    // Deploy testovací aplikaci (použijte sample-resources!)
    deployTestApp(t, suite, ns)

    // Počkat na ready
    err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
    suite.Assert.AssertNoError(err)

    // Spustit CLI příkaz
    result := suite.RunCLI("export", "--namespace", ns)

    // Validovat výsledek
    suite.Assert.AssertCommandSuccess(result)
}
```

### Použití sample-resources

**DŮLEŽITÉ**: Vždy používejte existující `sample-resources/` aplikace!

```go
// ✅ SPRÁVNĚ - použití existující aplikace
func deployTestApp(t *testing.T, suite *framework.TestSuite, namespace string) {
    cmd := exec.Command("kubectl", "apply",
        "-f", "sample-resources/hello-world/manifest.yaml",
        "-n", namespace)
    cmd.Run()
}

// ❌ ŠPATNĚ - vytváření nových resources
func deployTestApp(...) {
    manifest := `apiVersion: apps/v1...`  // Nedělat!
}
```

### Dostupné sample aplikace

- **hello-world** (`sample-resources/hello-world/`)
  - Deployment: `apache-hello`
  - Service: `apache-hello-service`
  - Labels: `app=apache-hello`

- **wordpress** (`sample-resources/wordpress/`)
  - Deployments: `wordpress`, `wordpress-mysql`
  - Services: `wordpress`, `wordpress-mysql`
  - Má validační skripty

## CI/CD Integrace

### GitHub Actions

E2E testy se spouští automaticky při:

- Pull requestech modifikujících kód
- Push do main větve
- Manuálním spuštění workflow

Workflow: `.github/workflows/e2e-tests.yml`

### Lokální testování s Kind

```bash
# Vytvořit Kind cluster
kind create cluster --name e2e-test

# Spustit testy
make test-e2e

# Cleanup
kind delete cluster --name e2e-test
```

## Příklady testů

### Test základního exportu

```go
func testBasicExport(t *testing.T, suite *framework.TestSuite) {
    ns := suite.CreateTestNamespace("export-basic")
    deployTestApp(t, suite, ns)

    err := suite.Cluster.WaitForDeployment(ns, "apache-hello", 120*time.Second)
    suite.Assert.AssertNoError(err)

    exportDir := suite.CreateTempDir("export-basic")
    result := suite.RunCLI("export", "--namespace", ns, "--export-dir", exportDir)

    suite.Assert.AssertCommandSuccess(result)
    suite.Assert.AssertDirExists(filepath.Join(exportDir, "resources", ns))
    suite.Assert.AssertMinYAMLFiles(filepath.Join(exportDir, "resources", ns), 1)
}
```

### Test s label selectorem

```go
result := suite.RunCLI(
    "export",
    "--namespace", ns,
    "--label-selector", "app=apache-hello",
    "--export-dir", exportDir,
)
suite.Assert.AssertCommandSuccess(result)
```

### Test chybových stavů

```go
result := suite.RunCLI("export", "--namespace", "nonexistent")
// V závislosti na implementaci může uspět nebo selhat
if result.Failed() {
    suite.Assert.AssertStderrContains(result, "not found")
}
```

## Best Practices

### ✅ Doporučení

1. **Izolace testů**: Používejte `suite.CreateTestNamespace()`
2. **Cleanup**: Vždy `defer suite.Cleanup()`
3. **Sample resources**: Používejte `sample-resources/` aplikace
4. **Assertions**: Používejte assertion helpers místo ručních kontrol
5. **Logování**: Logujte důležité kroky pomocí `suite.LogInfo()`
6. **Čekání**: Vždy čekejte na ready stav resources
7. **Short mode**: Označte dlouhé testy pomocí `suite.SkipIfShort()`

### ❌ Co nedělat

1. Nevytvářejte inline YAML manifesty
2. Nezapomínejte na cleanup
3. Nepřeskakujte `SkipIfShort()` pro cluster testy
4. Nepoužívejte raw error checking
5. Nekódujte natvrdo namespace názvy
6. Nepředpokládejte, že resources jsou okamžitě ready

## Troubleshooting

### Testy selhávají s "cluster not found"

```bash
kubectl cluster-info
# Ujistěte se, že cluster běží
```

### Testy timeoutují

```bash
# Zvyšte timeout
go test -v ./test/e2e/... -timeout 60m
```

### Binary nenalezen

```bash
make build
ls -la ./bin/kubectl-migrate
```

### Verbose output

```bash
go test -v ./test/e2e/... -test.v
```

## Přidání nových testů

1. Vytvořte `test/e2e/mycommand_test.go`
2. Naimportujte framework
3. Napište testy podle vzoru v `export_test.go`
4. Přidejte Makefile target
5. Aktualizujte GitHub Actions workflow (pokud potřeba)

Příklad Makefile targetu:

```makefile
test-e2e-mycommand: build ## Run mycommand E2E tests
	@echo "Running mycommand E2E tests..."
	@E2E_BINARY=$(BUILD_DIR)/$(BINARY_NAME) $(GOTEST) -v -timeout 15m ./test/e2e/mycommand_test.go ./test/e2e/framework/*.go
```

## Prostředí proměnné

- `E2E_BINARY` - Cesta k kubectl-migrate binárce (default: `./bin/kubectl-migrate`)
- `KUBECONFIG` - Cesta ke kubeconfig souboru (používá default pokud není nastaveno)

## Budoucí vylepšení

- [ ] Testy pro `convert` příkaz
- [ ] Testy pro `apply` příkaz
- [ ] Testy pro `transform` příkaz
- [ ] Testy pro `transfer-pvc` příkaz
- [ ] Multi-cluster testy
- [ ] Performance benchmarky

## Odkazy

- **Uživatelská dokumentace**: [test/e2e/README.md](../test/e2e/README.md)
- **Příspěvatelská dokumentace**: [test/e2e/CONTRIBUTING.md](../test/e2e/CONTRIBUTING.md)
- **Framework kód**: [test/e2e/framework/](../test/e2e/framework/)
- **Příklad testů**: [test/e2e/export_test.go](../test/e2e/export_test.go)

## Podpora

Pro otázky nebo problémy:

- Otevřete issue v GitHub repository
- Podívejte se na existující testy jako příklady
- Zkontrolujte framework kód v `test/e2e/framework/`
