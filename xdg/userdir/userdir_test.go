/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package userdir

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestParseValue(t *testing.T) {
	homeDir := "/home/test"
	Convey("parseValue", t, func(c C) {
		value, err := parseValue([]byte(`"$HOME/Desktop/"`), homeDir)
		c.So(err, ShouldBeNil)
		c.So(value, ShouldEqual, "/home/test/Desktop")

		value, err = parseValue([]byte(`"/home/test/DesktopA"`), homeDir)
		c.So(err, ShouldBeNil)
		c.So(value, ShouldEqual, "/home/test/DesktopA")

		value, err = parseValue([]byte(`"$HOME/"`), homeDir)
		c.So(err, ShouldBeNil)
		c.So(value, ShouldEqual, "/home/test")

		value, err = parseValue([]byte(`"/"`), homeDir)
		c.So(err, ShouldBeNil)
		c.So(value, ShouldEqual, "/")

		value, err = parseValue([]byte(""), homeDir)
		c.So(err, ShouldNotBeNil)
		c.So(value, ShouldEqual, "")

		value, err = parseValue([]byte("$HOME"), homeDir)
		c.So(err, ShouldNotBeNil)
		c.So(value, ShouldEqual, "")

		value, err = parseValue([]byte(`"not abs"`), homeDir)
		c.So(err, ShouldNotBeNil)
		c.So(value, ShouldEqual, "")

	})
}

func TestParseUserDirsConfig(t *testing.T) {
	Convey("parseUserDirsConfig", t, func(c C) {
		os.Setenv("HOME", "/home/test")
		cfg, err := parseUserDirsConfig("./testdata/user-dirs.dirs")
		c.So(err, ShouldBeNil)
		c.So(cfg, ShouldResemble, map[string]string{"XDG_DESKTOP_DIR": "/home/test/??????", "XDG_DOCUMENTS_DIR": "/home/test/??????", "XDG_DOWNLOAD_DIR": "/home/test/??????", "XDG_MUSIC_DIR": "/home/test/??????", "XDG_PICTURES_DIR": "/home/test/??????", "XDG_PUBLICSHARE_DIR": "/home/test/.Public", "XDG_TEMPLATES_DIR": "/home/test/.Templates", "XDG_VIDEOS_DIR": "/home/test/??????"})
	})
}

func TestGet(t *testing.T) {
	Convey("Get", t, func(c C) {
		os.Setenv("HOME", "/home/test")
		testDataDir, err := filepath.Abs("./testdata")
		c.So(err, ShouldBeNil)

		os.Setenv("XDG_CONFIG_HOME", testDataDir)

		c.So(Get(Desktop), ShouldEqual, "/home/test/??????")
		c.So(Get(Download), ShouldEqual, "/home/test/??????")
		c.So(Get(Templates), ShouldEqual, "/home/test/.Templates")
		c.So(Get(PublicShare), ShouldEqual, "/home/test/.Public")
		c.So(Get(Documents), ShouldEqual, "/home/test/??????")
		c.So(Get(Music), ShouldEqual, "/home/test/??????")
		c.So(Get(Pictures), ShouldEqual, "/home/test/??????")
		c.So(Get(Videos), ShouldEqual, "/home/test/??????")
		c.So(Get("XXXX"), ShouldEqual, "/home/test")
	})
}

func TestReloadCache(t *testing.T) {
	Convey("ReloadCache", t, func(c C) {
		os.Setenv("HOME", "/home/test")
		testDataDir, err := filepath.Abs("./testdata")
		c.So(err, ShouldBeNil)

		os.Setenv("XDG_CONFIG_HOME", testDataDir)
		c.So(Get(Desktop), ShouldEqual, "/home/test/??????")

		testDataDir2, err := filepath.Abs("./testdata2")
		c.So(err, ShouldBeNil)
		os.Setenv("XDG_CONFIG_HOME", testDataDir2)
		err = ReloadCache()
		c.So(err, ShouldBeNil)
		c.So(Get(Desktop), ShouldEqual, "/home/test/MyDesktop")
	})
}
