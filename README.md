# Setup
```
docker run -e "SERVERID=123456789" -p 8080:8080 docker.pkg.github.com/tischi86/discord-exporter/discord-exporter:latest
```
http://localhost:8080/metrics

Add a scrape config to prometheus
```
scrape_configs:
  - job_name: discord
    static_configs:
      - targets: ['localhost:8080']
```


# Example Output
```
# HELP discord_channel_users Represents the number of discord users by channel
# TYPE discord_channel_users gauge
discord_channel_users{channel="Chillout"} 1
# HELP discord_total_users Represents the number of total discord users
# TYPE discord_total_users gauge
discord_total_users 1
# HELP discord_users Represents the online users
# TYPE discord_users gauge
discord_users{user="Tischi"} 1
```
