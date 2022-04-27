# Kyverno Automate Performance Test

Kyverno is a Kubernetes native policy engine that secures and automates Kubernetes configurations. 
This project automates scalability tests for Kyverno on large Kubernetes clusters with several namespaces and resources


## Test scenario
This test scenario will create loads of Kubernetes objects (Pod, Namespace, Deployment, Cronjob, ConfigMap, Secret) based on user defined scale

Scales mapping:
```
  xs: 100 total resource
  small: 500 total resource
  medium: 1000 total resource
  large: 2000 total resource 
  xl: 3000 total resource
```

## Getting Started

```
git clone https://github.com/husnialhamdani/kyvernop.git
cd kyvernop
go build .
```
  
Start automation
```
./kyvernop execute --scale medium
``` 


Cleanup
```
./kyvernop cleanup -size 500
```

## Anomaly Detection

Isolation Forest is an algorithm that detects anomalies by taking a subset of data and constructing many isolation trees out of it.

An isolation tree is constructed by randomly selecting a feature and randomly selecting a value from that feature. A forest is constructed by aggregating all the isolation trees.

We pass the the Kyverno usage as input data and this algorithm will provide a prediction, The isolation forest assigns 0 to the anomalous data and 1 to the normal data and finally it plot the anomalies predicted by Isolation forest.


## Report

After the automation has completed, the tools will automatically generate a report based on Kyverno performance behaviour during the test and using the algorithm mentioned above.

![alt text](https://github.com/husnialhamdani/kyvernop/blob/main/report.png?raw=true)

## Kyverno Consideration on large cluster

[Scales]  [Expected behaviour]

### Recommendation setup
