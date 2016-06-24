package sample

import (
	"fmt"
	"os"
	"github.com/wxjs33/napi/config"
)

type SampleConfig struct {
	maddr   string /* mysql addr */
	dbname  string /* db name */
	dbuser  string /* db username */
	dbpwd   string /* db password */

	api_loc string /* sample api location */
	loc     string /* sample location */
}


func (conf *SampleConfig) ParseConfig(cf *config.Config) error {
	var err error
	conf.maddr, err = cf.C.GetString("sample", "mysql_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Sample] Read conf: No mysql_addr")
		return err
	}
	conf.dbname, err = cf.C.GetString("sample", "mysql_dbname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Sample] Read conf: No mysql_dbname")
		return err
	}
	conf.dbuser, err = cf.C.GetString("sample", "mysql_dbuser")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Sample] Read conf: No mysql_dbuser")
		return err
	}
	conf.dbpwd, err = cf.C.GetString("sample", "mysql_dbpwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Sample] Read conf: No mysql_dbpwd")
		return err
	}

	conf.loc, err = cf.C.GetString("sample", "location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Sample] Read conf: No sample_location, use default location", SAMPLE_LOCATION)
		conf.loc = SAMPLE_LOCATION
	}

	conf.api_loc, err = cf.C.GetString("sample", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Sample] Read conf: No api_location, use default location", SAMPLE_API_LOCATION)
		conf.api_loc = SAMPLE_API_LOCATION
	}

	return nil
}
