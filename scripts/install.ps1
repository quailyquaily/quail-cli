param(
    [string]$InstallDir = $(if ($env:QUAIL_CLI_INSTALL_DIR) { $env:QUAIL_CLI_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA "Programs\quail-cli\bin" }),
    [switch]$NoPathUpdate
)

$ErrorActionPreference = "Stop"

$Repo = "quailyquaily/quail-cli"
$BinaryName = "quail-cli"

switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { $Arch = "x86_64" }
    "ARM64" { $Arch = "arm64" }
    "x86" { $Arch = "i386" }
    default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
}

$Asset = "${BinaryName}_Windows_${Arch}.zip"
$Url = "https://github.com/${Repo}/releases/latest/download/${Asset}"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("quail-cli-install-" + [System.Guid]::NewGuid().ToString())
$ArchivePath = Join-Path $TempDir $Asset

try {
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    Write-Host "Downloading $Url"
    Invoke-WebRequest -Uri $Url -OutFile $ArchivePath

    Expand-Archive -Path $ArchivePath -DestinationPath $TempDir -Force

    $Binary = Get-ChildItem -Path $TempDir -Recurse -File -Filter "${BinaryName}.exe" | Select-Object -First 1
    if (-not $Binary) {
        throw "Binary not found in release archive"
    }

    $Target = Join-Path $InstallDir "${BinaryName}.exe"
    Copy-Item -Path $Binary.FullName -Destination $Target -Force

    Write-Host "Installed $BinaryName to $Target"
    & $Target version

    if (-not $NoPathUpdate) {
        $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $PathItems = @()
        if ($UserPath) {
            $PathItems = $UserPath -split ";" | Where-Object { $_ }
        }

        $AlreadyInPath = $false
        foreach ($Item in $PathItems) {
            if ($Item.TrimEnd("\") -ieq $InstallDir.TrimEnd("\")) {
                $AlreadyInPath = $true
                break
            }
        }

        if (-not $AlreadyInPath) {
            $NewPath = if ($UserPath) { "$UserPath;$InstallDir" } else { $InstallDir }
            [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
            $env:Path = "$env:Path;$InstallDir"
            Write-Host "Added $InstallDir to the user PATH. Open a new terminal if quail-cli is not found."
        }
    }
}
finally {
    if (Test-Path $TempDir) {
        Remove-Item -Path $TempDir -Recurse -Force
    }
}
