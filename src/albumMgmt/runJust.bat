cls
::go clean -i -r
::go build -o albumMgmt.exe albumMgmt.go

rm -r C:\Workspaces\My_Workspace\goProjects\src\github.com\ovace\albumMgmt\Pictures\dest4

albumMgmt.exe -v=true -cp=true -dr=false -src=C:\Workspaces\My_Workspace\goProjects\src\github.com\ovace\albumMgmt\Pictures\org -dest=C:\Workspaces\My_Workspace\goProjects\src\github.com\ovace\albumMgmt\Pictures\dest4 > ..\logs\out.log 2>&1