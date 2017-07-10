# snap-plugin-collector-kafka
A snap collector plugin that collect Kafka metrics from the MX4J http adaptor. This plugin should support any service exposed with MX4J, though not tested on other service yet.

## how to build:
    make 
    
   This command will perform 1. cleaning, 2. dependencies installation with glide, 3. plugin compilation with go, and 4. a simple test located in ./kafka folder. After these steps, the plugin should locate in **./build/snap-plugin-collector-kafka.**
   
## how to use 
      snaptel plugin load \<plugin binary filepath\>
    
## configuration:
  Remember to add these fields below to snapteld global configuration (i.e. snapteld --config \<config-file\>) to provide mx4j server information.
  ```
    plugins:
      all:
        mx4j_url: localhost
        mx4j_port: 8082
  ```
For detail format please see snapteld-config.yaml or [here](https://github.com/intelsdi-x/snap/blob/87826917b0685927135dc78d4d2355ac9cba3dbb/examples/configs/snap-config-sample.yaml)

## note:
  * Tested under **go version go1.8.1 linux/amd64**
  * The plugin has to be loaded **AFTER** the kafka service was launched. Because it will pull available metric types during runtime from mx4j api when it is loading into snaptel. 

## TODO:
  * Caching getMetricAPI and getMetricType
  
