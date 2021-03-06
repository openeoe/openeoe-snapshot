package snapshot

import (

	// "openeoe/openeoe/migration/pkg/apis"

	"context"
	"encoding/json"
	"fmt"
	nanumv1alpha1 "openeoe/openeoe/apis/snapshot/v1alpha1"
	"openeoe/openeoe/omcplog"
	"openeoe/openeoe/openeoe-snapshot/pkg/util"
	"openeoe/openeoe/openeoe-snapshot/pkg/util/etcd"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	// "openeoe/openeoe/migration/pkg/controller"
)

//volumeSnapshotRun 내에는 PV 만 들어온다고 가정한다.
func etcdSnapshotRun(r *reconciler, snapshotSource *nanumv1alpha1.SnapshotSource, groupSnapshotKey string) (string, error) {
	omcplog.V(4).Info(snapshotSource)

	omcplog.V(3).Info("etcd snapshot start")
	etcdSnapshotKeyAllPath := util.MakeSnapshotKeyForSnapshot(groupSnapshotKey, snapshotSource)
	omcplog.V(3).Info(etcdSnapshotKeyAllPath)
	//etcdSnapshotKeyAllPath := util.MakeSnapshotKeyAllPath(groupSnapshotKey, etcdSnapshotKey)
	//snapshotSource.ResourceSnapshotKey = etcdSnapshotKeyAllPath

	//Client 로 데이터 가져오기.
	resourceJSONString, err := GetResourceJSON(snapshotSource)
	if err != nil {
		omcplog.Error("etcdsnapshot.go : GetResourceJSON for cluster error")
		return etcdSnapshotKeyAllPath, err
	}

	omcplog.V(2).Info("Input ETCD")
	omcplog.V(2).Info("  key : " + etcdSnapshotKeyAllPath)
	//ETCD 에 삽입
	etcdCtl, etcdInitErr := etcd.InitEtcd()
	if etcdInitErr != nil {
		omcplog.Error("etcdsnapshot.go : Etcd Init Err")
		return etcdSnapshotKeyAllPath, etcdInitErr
	}
	_, etcdPutErr := etcdCtl.Put(etcdSnapshotKeyAllPath, resourceJSONString)
	if etcdPutErr != nil {
		omcplog.Error("etcdsnapshot.go : Etcd Put Err")
		return etcdSnapshotKeyAllPath, etcdPutErr
	}
	//snapshotSource.SnapshotKey = snapshotKey
	omcplog.V(2).Info("Input ETCD end")
	return etcdSnapshotKeyAllPath, nil
}

// GetResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func GetResourceJSON(snapshotSource *nanumv1alpha1.SnapshotSource) (string, error) {

	var resourceObj runtime.Object

	client := cm.Cluster_genClients[snapshotSource.ResourceCluster]

	switch snapshotSource.ResourceType {
	case util.DEPLOY:
		resourceObj = &appsv1.Deployment{}
	case util.SERVICE:
		resourceObj = &corev1.Service{}
	case util.PVC:
		resourceObj = &corev1.PersistentVolumeClaim{}
	case util.PV:
		resourceObj = &corev1.PersistentVolume{}
	default:
		omcplog.Error("Invalid resourceType")
		return "", fmt.Errorf("Invalid resourceType")
	}
	client.Get(context.TODO(), resourceObj, snapshotSource.ResourceNamespace, snapshotSource.ResourceName)

	omcplog.V(3).Info("resourceType : " + snapshotSource.ResourceType + ", resourceName : " + snapshotSource.ResourceName + ", resourceNamespace: " + snapshotSource.ResourceNamespace)
	ret, err := obj2JsonString(resourceObj)
	if err != nil {
		omcplog.Error("GetResourceJSON for cluster error")
		omcplog.V(2).Info("Json Convert Error")
	}
	return ret, nil
}

// Obj2JsonString : Deployment 등과 같은 interface 를 json string 으로 변환.
func obj2JsonString(obj interface{}) (string, error) {

	json, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	omcplog.V(3).Info("===Obj2JsonString===")
	omcplog.V(3).Info(string(json)[0:40] + "...")

	return string(json), nil
}
