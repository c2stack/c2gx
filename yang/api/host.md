

# OS Metrics


## <a name=""></a>/
OS metrics for this host


  
* **[swap](#/swap)** - . 

  
* **[ram](#/ram)** - OS memory details. 

  
* **[fs[…]](#/fs)** - file system details. 

  
* **[cpu[…]](#/cpu)** - current workload for each CPU. 

  
* **[proc[…]](#/proc)** - . 







## <a name="/swap"></a>/swap/



  
* **total** `uint64` - . 

  
* **used** `uint64` - . 

  
* **free** `uint64` - . 







## <a name="/ram"></a>/ram/
OS memory details


  
* **total** `uint64` - . 

  
* **used** `uint64` - . 

  
* **free** `uint64` - . 

  
* **actualFree** `uint64` - . 

  
* **actualUsed** `uint64` - . 







## <a name="/fs"></a>/fs={dirName}/
file system details


  
* **dirName** `string` - . 

  
* **devName** `string` - . 

  
* **typeName** `string` - . 

  
* **sysTypeName** `string` - . 

  
* **options** `string` - . 

  
* **flags** `int32` - . 

  
* **[usage](#/fs/usage)** - file system usage. 







## <a name="/fs/usage"></a>/fs={dirName}/usage/
file system usage


  
* **total** `uint64` - . 

  
* **used** `uint64` - . 

  
* **free** `uint64` - . 

  
* **avail** `uint64` - . 

  
* **files** `uint64` - . 

  
* **freeFiles** `uint64` - . 







## <a name="/cpu"></a>/cpu={id}/
current workload for each CPU


  
* **id** `int32` - . 

  
* **user** `uint64` - . 

  
* **nice** `uint64` - . 

  
* **sys** `uint64` - . 

  
* **idle** `uint64` - . 

  
* **wait** `uint64` - . 

  
* **irq** `uint64` - . 

  
* **softIrq** `uint64` - . 

  
* **stolen** `uint64` - . 







## <a name="/proc"></a>/proc={pid}/



  
* **pid** `int32` - . 

  
* **[state](#/proc/state)** - . 

  
* **[mem](#/proc/mem)** - . 

  
* **[time](#/proc/time)** - . 







## <a name="/proc/state"></a>/proc={pid}/state/



  
* **name** `string` - . 

  
* **state** `enumeration` - .  *Allowed Values: sleep,run,stop,zombie,idle,unknown* 

  
* **ppid** `int32` - Parent PID. 

  
* **tty** `int32` - . 

  
* **priority** `int32` - . 

  
* **nice** `int32` - . 

  
* **processor** `int32` - . 







## <a name="/proc/mem"></a>/proc={pid}/mem/



  
* **size** `uint64` - . 

  
* **resident** `uint64` - . 

  
* **share** `uint64` - . 

  
* **minorFaults** `uint64` - . 

  
* **majorFaults** `uint64` - . 

  
* **pageFaults** `uint64` - . 







## <a name="/proc/time"></a>/proc={pid}/time/



  
* **startTime** `uint64` - . 

  
* **user** `uint64` - . 

  
* **sys** `uint64` - . 

  
* **total** `uint64` - . 







