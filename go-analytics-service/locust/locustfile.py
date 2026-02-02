import random
import time
from locust import HttpUser, task, between
import json

class MetricsUser(HttpUser):
    wait_time = between(0.01, 0.1)  # 10-100ms between requests (10-100 RPS per user)
    
    def on_start(self):
        """Called when a simulated user starts"""
        self.base_rps = 100.0
        self.cpu_base = 50.0
        
    @task(10)
    def send_normal_metric(self):
        """Send normal metrics (90% of traffic)"""
        timestamp = int(time.time())
        metric = {
            "timestamp": timestamp,
            "cpu": self.cpu_base + random.gauss(0, 5),  # Normal distribution
            "rps": self.base_rps + random.gauss(0, 10)
        }
        self.client.post(
            "/metrics",
            json=metric,
            name="/metrics",
            catch_response=True
        )
    
    @task(1)
    def send_anomalous_metric(self):
        """Send anomalous metrics (10% of traffic)"""
        timestamp = int(time.time())
        # Create anomaly: high RPS spike
        metric = {
            "timestamp": timestamp,
            "cpu": self.cpu_base + random.gauss(0, 5),
            "rps": self.base_rps + random.gauss(150, 30)  # High spike
        }
        self.client.post(
            "/metrics",
            json=metric,
            name="/metrics",
            catch_response=True
        )
    
    @task(1)
    def get_analytics(self):
        """Query analytics endpoint"""
        self.client.get("/analyze", name="/analyze")
    
    @task(1)
    def health_check(self):
        """Health check"""
        self.client.get("/health", name="/health")

# Configuration for 1000 RPS
# Calculate: 1000 RPS / ~50 RPS per user = ~20 users
# With wait_time of 0.01-0.1, each user generates 10-100 RPS
# So we need ~20 users for 1000 RPS
