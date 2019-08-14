// Copyright 2018-2019 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package main

import (
	"os"
	"time"

	rpcpb "github.com/cs3org/go-cs3apis/cs3/rpc"
	usershareproviderv0alphapb "github.com/cs3org/go-cs3apis/cs3/usershareprovider/v0alpha"
	"github.com/jedib0t/go-pretty/table"
)

func shareListCommand() *command {
	cmd := newCommand("share-list")
	cmd.Description = func() string { return "list shares you manage" }
	cmd.Usage = func() string { return "Usage: share list [-flags]" }
	cmd.Action = func() error {
		ctx := getAuthContext()
		shareClient, err := getUserShareProviderClient()
		if err != nil {
			return err
		}

		shareRequest := &usershareproviderv0alphapb.ListSharesRequest{}

		shareRes, err := shareClient.ListShares(ctx, shareRequest)
		if err != nil {
			return err
		}

		if shareRes.Status.Code != rpcpb.Code_CODE_OK {
			return formatError(shareRes.Status)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Owner.Idp", "Owner.OpaqueId", "ResourceId", "Permissions", "Type", "Grantee.Idp", "Grantee.OpaqueId", "Created", "Updated"})

		for _, s := range shareRes.Shares {
			t.AppendRows([]table.Row{
				{s.Id.OpaqueId, s.Owner.Idp, s.Owner.OpaqueId, s.ResourceId.String(), s.Permissions.String(), s.Grantee.Type.String(), s.Grantee.Id.Idp, s.Grantee.Id.OpaqueId, time.Unix(int64(s.Ctime.Seconds), 0), time.Unix(int64(s.Mtime.Seconds), 0)},
			})
		}
		t.Render()
		return nil
	}
	return cmd
}