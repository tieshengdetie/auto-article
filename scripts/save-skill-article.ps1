[CmdletBinding()]
param(
    [Parameter(Mandatory = $true, Position = 0)]
    [Alias("Path")]
    [string]$Payload,

    [string]$BaseUrl = $env:AUTO_ARTICLE_BASE_URL,

    [int]$Timeout = 30,

    [switch]$DryRun,

    [string[]]$StaticRoot = @(),

    [Parameter(ValueFromPipeline = $true)]
    [string]$PipelineInput
)

begin {
    $ErrorActionPreference = "Stop"
    $stdinBuffer = [System.Collections.Generic.List[string]]::new()
}

process {
    if ($null -ne $PipelineInput) {
        $stdinBuffer.Add($PipelineInput)
    }
}

end {
    $repoRoot = (Resolve-Path -LiteralPath (Join-Path $PSScriptRoot "..")).Path
    $saveScript = Join-Path $repoRoot "skills\auto-media-writer\scripts\save_skill_article.py"

    if (-not (Test-Path -LiteralPath $saveScript)) {
        Write-Error "Cannot find auto-media-writer save script: $saveScript"
        exit 2
    }

    $python = if ($env:PYTHON) {
        $env:PYTHON
    } elseif (Get-Command python3 -ErrorAction SilentlyContinue) {
        "python3"
    } else {
        "python"
    }
    $argsList = @($saveScript, $Payload, "--timeout", "$Timeout")

    if ($BaseUrl) {
        $argsList += @("--base-url", $BaseUrl)
    }

    if ($DryRun) {
        $argsList += "--dry-run"
    }

    foreach ($root in $StaticRoot) {
        if ($root) {
            $argsList += @("--static-root", $root)
        }
    }

    if ($Payload -eq "-" -and $stdinBuffer.Count -gt 0) {
        $stdinText = [string]::Join([Environment]::NewLine, $stdinBuffer)
        $stdinText | & $python @argsList
    } else {
        & $python @argsList
    }
    exit $LASTEXITCODE
}
