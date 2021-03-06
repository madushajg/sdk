/*
 * Copyright (c) 2019 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package gcp

import (
	"fmt"
	"strings"
	"time"
)

func (gcp *Gcp) TearDown() error {
	uniqueNumber := strings.TrimPrefix(gcp.clusterName, "cellery-cluster")
	_, err := gcp.service.Projects.Zones.Clusters.Delete(gcp.projectName, gcp.zone, clusterNamePrefix+uniqueNumber).Do()
	if err != nil {
		return fmt.Errorf("failed to delete gcp cluster: %v", err)

	}
	for i := 0; i < 15; i++ {
		if gcp.gcpClusterExist() {
			time.Sleep(60 * time.Second)
		} else {
			break
		}
	}
	_, err = gcp.sqlService.Instances.Delete(gcp.projectName, dbInstanceNamePrefix+uniqueNumber).Do()
	if err != nil {
		return fmt.Errorf("failed to delete the sql instance: %v", err)
	}
	_, err = gcp.nfsService.Projects.Locations.Instances.Delete("projects/" + gcp.projectName + "/locations/" + gcp.zone + "/instances/" + fileStorePrefix + uniqueNumber).Do()
	if err != nil {
		return fmt.Errorf("failed to delete nfs server %v", err)
	}
	object := gcp.storageClient.Bucket(storagePrefix + uniqueNumber).Object(initSql)
	if err := object.Delete(gcp.ctx); err != nil {
		return fmt.Errorf("error deleting gcp storage object: %v", err)
	}
	if err := gcp.storageClient.Bucket(storagePrefix + uniqueNumber).Delete(gcp.ctx); err != nil {
		return err
	}
	return nil
}
