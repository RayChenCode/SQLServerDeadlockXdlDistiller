# SQLServerDeadlockXdlDistiller
Parsing SQL server deadlock xdl files to CSV format!

### Installation & go

- Generate go binary execute file.
```sh
go build .
```
- Put all deadlock xdl files into a directory.
- Double click execute file.
- You get a deadlock_report.csv!

### Configuration

Setting in config.yaml

```sh
 deadlock:
     directory:
         # Setting the directory of deadlock xdl files at here
         path: C:\\Users\\ray.chen\\Desktop\\sqldeadlock20180802
     report:
         # Setting the file of output report csv file at here
         filepath: deadlock_report.csv
```
### Deadlock report information

- 1. The date of deadlock occured
- 2. The datetime of deadlock occured
- 3. Deadlock trigged table
- 4. XDL filename
- 5. Victim SQL
