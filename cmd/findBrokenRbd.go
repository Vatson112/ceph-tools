/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/Vatson112/ceph-tools/internal"
	"github.com/ceph/go-ceph/rbd"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

// findBrokenRbdCmd represents the findBrokenRbd command
var findBrokenRbdCmd = &cobra.Command{
	Use:   "findBrokenRbd",
	Short: "Find incorrectly deleted rbd images",
	Long: `Find incorrectly deleted rbd images.
When image deleted from cinder and something happen with cinder at this time rbd image may be in incorrect state.
It is showed in list of images, but we cannot open it.
This command find these images by getting list of images in pull and trying to open it. When it failed - image is incorrect.

See ref:
- https://tracker.ceph.com/issues/52910
- https://www.spinics.net/lists/ceph-users/msg75810.html
- https://lists.ceph.io/hyperkitty/list/ceph-users@ceph.io/thread/5RQRBZOFPBVR77QHH6PV2UR4KQSGR44U/
-

You must set --poolname for using it.`,
	Run: func(cmd *cobra.Command, args []string) {
		cconf, _ := cmd.Flags().GetString("ceph-config")
		pname, _ := cmd.Flags().GetString("poolname")
		delete, _ := cmd.Flags().GetBool("delete")
		rcon := internal.NewRadosConnection(cconf)
		defer rcon.Shutdown()
		err := rcon.Connect()
		if err != nil {
			slog.Error(fmt.Sprintf("Error login to cluster. Err=%s", err))
		}
		fsid, _ := rcon.Client.GetFSID()
		slog.Info(fmt.Sprintf("Login to Ceph cluster=%s", fsid))
		ioctx, err := rcon.Client.OpenIOContext(pname)
		if err != nil {
			slog.Error(fmt.Sprintf("Error opening pool. Err=%s", err))
		}
		defer ioctx.Destroy()
		images, err := rbd.GetImageNames(ioctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Error getting image names. Err=%s", err))
		}
		slog.Info(fmt.Sprintf("Using pool=%s", pname))
		brokenImage := make([]string, 0, 0)
		for _, iname := range images {
			image := rbd.GetImage(ioctx, iname)
			slog.Debug(fmt.Sprintf("Find image=%s", image.GetName()))
			err := image.Open()
			if err != nil {
				brokenImage = append(brokenImage, iname)
			}
			// metadata, err := image.Stat()
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// pp.Println(metadata)
		}
		slog.Info("Found broken images:")
		pp.Println(brokenImage)
		if delete {
			slog.Info("Start removing broken images.")
			for _, iname := range brokenImage {
				image := rbd.GetImage(ioctx, iname)
				slog.Debug(fmt.Sprintf("Delete image=%s", image.GetName()))
				err := image.Remove()
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to delete image=%s", iname))
				}
				slog.Info(fmt.Sprintf("Deleted image=%s", image.GetName()))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(findBrokenRbdCmd)
	findBrokenRbdCmd.Flags().String("poolname", "", "Name of the rados pool")
	findBrokenRbdCmd.Flags().Bool("delete", false, "Delete founded rbd images.")
	findBrokenRbdCmd.MarkFlagRequired("poolname")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findBrokenRbdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findBrokenRbdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type findBrokenRbd struct {
	poolName string
}

func new(pool string) *findBrokenRbd {
	return &findBrokenRbd{
		poolName: pool,
	}
}
func (findBrokenRbd) run() {

}
