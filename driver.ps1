function Run-SingleTest ($testPath) {
    go run main.go $testPath
}

function Run-AllTests {
    $testFiles = Get-ChildItem -Path .\probs
    $testFiles | ForEach-Object { Run-SingleTest $_.FullName }
}

Run-AllTests