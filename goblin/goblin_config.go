package goblin

import (
	"fmt"
	"os"
	"github.com/wxjs33/napi/config"
)

type GoblinConfig struct {
	maddr      string  /* mysql addr */
	dbname     string  /* db name */
	dbuser     string  /* db username */
	dbpwd      string  /* db password */

	api_loc    string  /* goblin api location */
	goblin_loc string  /* goblin location */
}


func (conf *GoblinConfig) ParseConfig(cf *config.Config) error {
	var err error
	conf.maddr, err = cf.C.GetString("goblin", "mysql_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_addr")
		return err
	}
	conf.dbname, err = cf.C.GetString("goblin", "mysql_dbname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbname")
		return err
	}
	conf.dbuser, err = cf.C.GetString("goblin", "mysql_dbuser")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbuser")
		return err
	}
	conf.dbpwd, err = cf.C.GetString("goblin", "mysql_dbpwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [Goblin] Read conf: No mysql_dbpwd")
		return err
	}

	conf.goblin_loc, err = cf.C.GetString("goblin", "goblin_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Goblin] Read conf: No goblin_location, use default location", DEFAULT_GOBLIN_LOCATION)
		conf.goblin_loc = DEFAULT_GOBLIN_LOCATION
	}

	conf.api_loc, err = cf.C.GetString("goblin", "api_location")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Goblin] Read conf: No api_location, use default location", DEFAULT_API_LOCATION)
		conf.api_loc = DEFAULT_API_LOCATION
	}

	return nil
}
