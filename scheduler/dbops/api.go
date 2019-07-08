package dbops

import "log"

//user ->api service ->delete video
//api service->scheduler->write video deletion record
//timer->runner->read video deletion record->executor->delete video from folder

func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("insert into video_del_rec(video_id)values ?")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("addvideodeleterecord error: %v", err)
		return err
	}
	defer stmtIns.Close()
	return nil
}
