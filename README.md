# Kyverno Automate Performance Test

Kyverno is a Kubernetes native policy engine that secures and automates Kubernetes configurations. 
This project automates scalability tests for Kyverno on large Kubernetes clusters with several namespaces and resources


## Test scenario



## Getting Started

```
  git clone https://github.com/husnialhamdani/kyvernop.git
  cd kyvernop
  go build .
```
  
Start automation
  ./kyvernop execute --scale medium
  
Cleanup
  ./kyvernop cleanup -size 500


## Anomaly Detection



## Report

After the automation has completed, the tools will automatically generate a report based on Kyverno performance behaviour during the test and using the algorithm mentioned above.

![alt text](https://github.com/husnialhamdani/kyvernop/blob/main/report.png?raw=true)
