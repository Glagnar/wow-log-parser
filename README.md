# GENISIS Log Parser

Output help and defaults
```bash
./wow-log-parser -help
```

Check inputfile
```bash
./wow-log-parser -checkonly=true -inputfile=wowlogs.txt
```

Run sort
```bash
./wow-log-parser -checkonly=false -inputfile=wowlogs.txt -outputfile=fixedoutput.txt
```

This works too, will use output.txt and input.txt defaults
```bash
./wow-log-parser -checkonly=false
```

This will split a single file, in the individual raids that are in it
```bash
./wow-log-parser -inputfile WoWCombatLog.txt -checkonly=false -split=true
```

If you want to output errors to screen add loglevel
```bash
./wow-log-parser -checkonly=true -inputfile=wowlogs.txt -loglevel=debug
```