# Rate Reader

## Run :

    go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"
       -config-path=./Listeners/Btc/config/config.yml for debug in JetBrains Goland
       
###  Docker :
    
    Build container:
    
    docker build --rm -t rate-reader:latest .
    
    Run container 
    docker run --rm -d --env-file=config/dev.env -p 5001:5001 --name rate-reader rate-reader:latest
    
    Read logs from container
    docker logs rate-reader