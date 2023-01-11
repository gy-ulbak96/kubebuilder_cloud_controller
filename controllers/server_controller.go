/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudcontrollerv1alpha1 "github.com/gy-ulbak96/kubebuilder_cloud_controller/api/v1alpha1"
)

const (
	cloudUrl = "http://127.0.0.1:8080"
)

// ServerReconciler reconciles a Server object
type ServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cloudcontroller.crd.example.com,resources=servers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloudcontroller.crd.example.com,resources=servers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloudcontroller.crd.example.com,resources=servers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Server object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *ServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger = log.FromContext(ctx).WithValues("req.Namespace", req.Namespace, "req.Name", req.Name)
	reqLogger.Info("Reconciling Servers.")

  //Check if the Server already exist, if not create a new one.
	server := &cloudcontrollerv1alpha1.Server{}
	err := r.Client.Get(context.TODO(), rq.NamespacedName, server)
	if err != nil{
		if errors.IsNotFound(err) {
			reqLogger.Info("Server resource not found.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Server.")
		return ctrl.Result{}, err
	}
	// err := r.Client.Get(context.TODO(), req.)
  
	serverName := server.Spec.ServerName
	if serverName == "" {
		utilruntime.HandleError(fmt.Errorf("%s: servername must be specified", key))
		return nil
	}

	client, err := cloudclient.CreateClient(cloudurl)
  {                                                                                                          
      utilruntime.HandleError(fmt.Errorf("%+v", err.Error()))
      return nil
  }

	// err = r.Client.Create(context.TODO(), dep)
	// if err != nil{
	// 	reqLogger.Error(err, "Failed to create new Server.")
	// 	return ctrl.Result{}, err
	// }

	if server.Status.ServerId == "" {
		deletionTimestamp := server.ObjectMeta.GetDeletionTimestamp()
		finalizer := server.ObjectMeta.GetFinalizers()

		if len(finalizer) > 0 && deletionTimestamp != nil {
			serverId := server.Status.ServerId
			_, errGetServer := client.GetServer(serverId)
			if errGetServer != nil {
                utilruntime.HandleError(fmt.Errorf("%+v", errGetServer.Error()))
                return nil
      }

			errDeleteServer := deleteServer(serverId)
            if errDeleteServer != nil {
                utilruntime.HandleError(fmt.Errorf("%+v", errDeleteServer.Error()))
                return nil
            }
            fmt.Println("Success to delete on " + serverName)
						// Server  object에서 finalizer를 제거합니다. 기존 server object를 직접 수정하지 않고 DeepCopy를 통해 복사한 object를 수정한 뒤 업데이트합니다.

            serverCopy := server.DeepCopy()
            serverCopy.ObjectMeta.SetFinalizers([]string{})
            _, errUpdate := cloudcontrollerv1alpha1.Servers(server.Namespace).Update(context.TODO(), serverCopy, metav1.UpdateOptions{})
            if errUpdate != nil {
                utilruntime.HandleError(fmt.Errorf("%+v", errUpdate.Error()))
                return nil
            }
            fmt.Println("Success to remove a finalizer on " + serverName)
            return nil
        }
		}
			return ctrl.Result{}, nil
}


	
	


func deleteServer(){}

func updateServerStatus(){}

// SetupWithManager sets up the controller with the Manager.
func (r *ServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudcontrollerv1alpha1.Server{}).
		Complete(r)
}
