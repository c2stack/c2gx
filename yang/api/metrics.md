

# Metrics Store


## <a name=""></a>/



  
* **addr** `string` - http address of influx db service. 

  
* **[relay[…]](#/relay)** - . 







## <a name="/relay"></a>/relay={name}/



  
* **name** `string` - . 

  
* **database** `string` - . 

  
* **[tag[…]](#/relay/tag)** - . 

  
* **script** `string` - . 

  
* **[source](#/relay/source)** - . 





### Events:

* <a name="/relay/update"></a>**/relay={name}/update** - 

 	
> * **time** `int64` - 	
> * **name** `string` - 	
> * **database** `string` - 
> * **tag[…]** - 
>     * **name** -  
>     * **value** -  
> * **field[…]** - 
>     * **name** -  
>     * **value** -  





## <a name="/relay/tag"></a>/relay={name}/tag={name}/



  
* **name** `string` - . 

  
* **value** `string` - . 







## <a name="/relay/source"></a>/relay={name}/source/



  
* **device** `string` - . 

  
* **module** `string` - . 

  
* **path** `string` - . 







