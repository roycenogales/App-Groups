# App-Groups

![Logo](images/Logo.png)

## What is App Groups?

App Groups is a Windows desktop application that allows users to save and restore groups of open applications and their window layouts.

It captures the state of selected application windows, including positioning and dimensions. A workspace can be named and saved, then to be restored at any time, automatically reopening applications in the same state as they were when the workspace was created.

![Application](images/Application.png)

## Features
- Window Detection: Detects all currently open windows.
- Window Selection: Allows users to choose which windows to save to a workspace.
- Persistent Storage: Stores workspaces locally in JSON.
- Workspace Restoration: Reopens applications and restores their saved window positions, dimensions, and states.
- Workspace Management: Allows users to view and delete saved workspaces.

## Technology
  App Groups is built using Go and Wails v3, with a JavaScript-based frontend.
  
- Go: Backend logic, window detection, and workspace persistence.
- Wails v3: Desktop application framework.
- JavaScript: Frontend interface and user interaction.
- JSON: Workspace configuration and persistence.

## Requirements

| 	|Version|
| -  	| -	|
| Windows| 10/11|
| Go 	| 1.26.0|
| Wails | 3.0.0 |
|Node.js|	|
| npm	|	|
|WebView2|	|

## Installing Dependencies

### Setup Powershell

[How to Setup Powershell for Github](https://github.com/roycenogales/SetUp-Powershell)

### Other Dependencies

```
winget install --id GoLang.Go -e
winget install --id OpenJS.NodeJS.LTS -e
winget install --id Git.Git -e
```

Restart PowerShell.

```
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
#This may take a while
```

Ensure %USERPROFILE%\go\bin is included in your PATH.

```
$env:Path -split ';' | Where-Object { $_ -like '*\go\bin' }
```

If nothing appears, run the next command.

```
$goBin = "$env:USERPROFILE\go\bin"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")

if (($userPath -split ';') -notcontains $goBin) {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$userPath;$goBin",
        "User"
    )
}
```

Restart Powershell and check the path again.

Verify the environment.

```
go version
node --version
npm --version
git --version
wails3 version
wails3 doctor
```

If npm is giving a restriction error, follow these steps:

```
#Check if CurrentUser or Local Machine is Restricted of Undefined.
Get-ExecutionPolicy -List

#Then change the policy.
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned

#Restart PowerShell
```

## Build from source

```
#Clone the Repository:

git clone https://github.com/roycenogales/App-Groups.git

#Navigate to the project directory:

cd App-Groups

#Build the application:

wails3 build

#the compiled application will be available in 'build/bin/'.
```

## Known Limitations

## Contributors
Royce Nogales [github.com/roycenogales]

### License
MIT

### References

