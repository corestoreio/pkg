// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package eav

import (
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
)

// GetAttributeSelectSql generates the select query to retrieve full attribute configuration
// Implements the scope on a SQL query basis so that attribute functions does not need to deal with it.
// Tests see the tools package
// @see magento2/app/code/Magento/Eav/Model/Resource/Attribute/Collection.php::_initSelect()
func GetAttributeSelectSql(dbrSess dbr.SessionRunner, aat EntityTypeAdditionalAttributeTabler, entityTypeID, websiteID int64) (*dbr.SelectBuilder, error) {

	ta, err := TableCollection.Structure(TableIndexAttribute)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("eav.GetAttributeSelectSql.TableCollection.Structure", "err", err, "entityTypeID", entityTypeID, "websiteID", websiteID)
		}
		return nil, errgo.Mask(err)
	}
	taa, err := aat.TableAdditionalAttribute()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("eav.GetAttributeSelectSql.TableAdditionalAttribute", "err", err, "ta", ta, "entityTypeID", entityTypeID, "websiteID", websiteID)
		}
		return nil, errgo.Mask(err)
	}
	tew, err := aat.TableEavWebsite()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("eav.GetAttributeSelectSql.TableEavWebsite", "err", err, "ta", ta, "taa", taa, "entityTypeID", entityTypeID, "websiteID", websiteID)
		}
		return nil, errgo.Mask(err)
	}
	// tew table can now contains columns names which can occur in table eav_attribute and
	// or [catalog|customer|entity]_eav_attribute
	var (
		ifNull           []string
		tewAddedCols     []string
		taColumnsQuoted  = util.StringSlice(ta.AllColumnAliasQuote(csdb.MainTable))
		taaColumnsQuoted = util.StringSlice(taa.ColumnAliasQuote(csdb.AdditionalTable))
	)

	if tew != nil {
		ifNull = make([]string, len(tew.Columns.ColumnsNoPK().FieldNames()))
		for i, tewC := range tew.Columns.ColumnsNoPK().FieldNames() {
			t := ""
			switch {
			case ta.In(tewC):
				t = csdb.MainTable
				break
			case taa.In(tewC):
				t = csdb.AdditionalTable
				break
			default:
				err := errgo.Newf("Cannot find column name %s.%s neither in table %s nor in %s.", tew.Name, tewC, ta.Name, taa.Name)
				if PkgLog.IsDebug() {
					PkgLog.Debug("eav.GetAttributeSelectSql.Columns.FieldNames.default", "err", err, "ta", ta, "taa", taa, "entityTypeID", entityTypeID, "websiteID", websiteID)
				}
				return nil, err
			}
			ifNull[i] = dbr.IfNullAs(csdb.ScopeTable, tewC, t, tewC, tewC)
			tewAddedCols = append(tewAddedCols, tewC)
		}
		taColumnsQuoted.ReduceContains(tewAddedCols...)
		taaColumnsQuoted.ReduceContains(tewAddedCols...)
	}

	selectSql := dbrSess.
		Select(taColumnsQuoted...).
		From(ta.Name, csdb.MainTable).
		Join(
			dbr.JoinTable(taa.Name, csdb.AdditionalTable),
			taaColumnsQuoted,
			dbr.ConditionRaw(csdb.AdditionalTable+".attribute_id = "+csdb.MainTable+".attribute_id"),
			dbr.ConditionRaw(csdb.MainTable+".entity_type_id = ?", entityTypeID),
		)

	if len(tewAddedCols) > 0 {
		selectSql.
			LeftJoin(
				dbr.JoinTable(tew.Name, csdb.ScopeTable),
				ifNull,
				dbr.ConditionRaw(csdb.ScopeTable+"."+"attribute_id"+" = "+csdb.MainTable+"."+"attribute_id"),
				dbr.ConditionRaw(csdb.ScopeTable+"."+"website_id"+" = ?", websiteID),
			)
	}
	return selectSql, nil
}
