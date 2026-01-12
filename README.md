**Build the Executable:**
(Run this once to create the `twin.exe` file)

```powershell
go build -o twin.exe .
```

**Start Temporal Server (if not already running in Docker):**

```bash
temporal server start-dev
```

**Start Twin A (Terminal 1):**

```powershell
.\twin.exe -twin=A -port1=8001 -port2=8002
```

**Start Twin B (Terminal 2):**

```powershell
.\twin.exe -twin=B -port1=8002 -port2=8001
```
