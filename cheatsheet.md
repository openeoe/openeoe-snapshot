

```
#delete openeoe snapshot resources.
kubectl delete pv,pvc,job  --selector=openeoe=snapshot -n openeoe  --context=cluster1



```