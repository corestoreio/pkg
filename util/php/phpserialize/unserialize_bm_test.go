package phpserialize_test

import (
	"github.com/corestoreio/csfw/util/php/phpserialize"
	"testing"
)

var benchmarkUnserialize phpserialize.PhpValue

// BenchmarkUnserialize-4	   10000	    174507 ns/op	   46955 B/op	    1408 allocs/op
func BenchmarkUnserialize(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var err error
		decoder := phpserialize.NewUnSerializer(toUnserialize)
		benchmarkUnserialize, err = decoder.Decode()
		if err != nil {
			b.Error(err)
		}
	}
}

var benchmarkSerialize string

// BenchmarkSerialize-4  	   10000	    126599 ns/op	   56254 B/op	     959 allocs/op
func BenchmarkSerialize(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var err error
		enc := phpserialize.NewSerializer()
		benchmarkSerialize, err = enc.Encode(toSerialize)
		if err != nil {
			b.Error(err)
		}
	}
}

const toUnserialize = `a:25:{i:0;a:4:{s:5:"order";s:2:"10";s:6:"column";s:2:"id";s:9:"attribute";s:21:"xunseria_directive_id";s:5:"param";s:1:"0";}i:1;a:3:{s:5:"order";s:2:"10";s:6:"column";s:13:"item_group_id";s:9:"attribute";s:32:"xunseria_directive_item_group_id";}i:2;a:4:{s:5:"order";s:2:"20";s:6:"column";s:5:"title";s:9:"attribute";s:30:"xunseria_directive_concatenate";s:5:"param";s:77:"attrib {{attrib_name_map}} {{attrib_serie}} {{name}} {{attrib_articlenumber}}";}i:3;a:3:{s:5:"order";s:2:"30";s:6:"column";s:11:"description";s:9:"attribute";s:11:"description";}i:4;a:4:{s:5:"order";s:2:"40";s:6:"column";s:4:"link";s:9:"attribute";s:22:"xunseria_directive_url";s:5:"param";s:27:"?utm_source=google_shopping";}i:5;a:4:{s:5:"order";s:2:"50";s:6:"column";s:10:"image_link";s:9:"attribute";s:29:"xunseria_directive_image_link";s:5:"param";s:5:"image";}i:6;a:4:{s:5:"order";s:2:"60";s:6:"column";s:21:"additional_image_link";s:9:"attribute";s:40:"xunseria_directive_additional_image_link";s:5:"param";s:5:"image";}i:7;a:4:{s:5:"order";s:2:"70";s:6:"column";s:5:"price";s:9:"attribute";s:24:"xunseria_directive_price";s:5:"param";s:1:"1";}i:8;a:4:{s:5:"order";s:2:"80";s:6:"column";s:10:"sale_price";s:9:"attribute";s:29:"xunseria_directive_sale_price";s:5:"param";s:1:"1";}i:9;a:3:{s:5:"order";s:2:"90";s:6:"column";s:25:"sale_price_effective_date";s:9:"attribute";s:44:"xunseria_directive_sale_price_effective_date";}i:10;a:3:{s:5:"order";s:3:"110";s:6:"column";s:12:"availability";s:9:"attribute";s:31:"xunseria_directive_availability";}i:12;a:4:{s:5:"order";s:3:"140";s:6:"column";s:5:"brand";s:9:"attribute";s:31:"xunseria_directive_static_value";s:5:"param";s:6:"attrib";}i:13;a:3:{s:5:"order";s:3:"150";s:6:"column";s:3:"mpn";s:9:"attribute";s:20:"attrib_articlenumber";}i:14;a:4:{s:5:"order";s:3:"160";s:6:"column";s:9:"condition";s:9:"attribute";s:31:"xunseria_directive_static_value";s:5:"param";s:3:"new";}i:15;a:3:{s:5:"order";s:3:"170";s:6:"column";s:12:"product_type";s:9:"attribute";s:43:"xunseria_directive_product_type_by_category";}i:16;a:3:{s:5:"order";s:3:"180";s:6:"column";s:23:"google_product_category";s:9:"attribute";s:46:"xunseria_directive_google_category_by_category";}i:17;a:3:{s:5:"order";s:3:"190";s:6:"column";s:17:"identifier_exists";s:9:"attribute";s:36:"xunseria_directive_identifier_exists";}i:18;a:3:{s:5:"order";s:3:"200";s:6:"column";s:9:"is_bundle";s:9:"attribute";s:28:"xunseria_directive_is_bundle";}s:18:"_1447682911167_167";a:3:{s:5:"order";s:3:"210";s:6:"column";s:5:"color";s:9:"attribute";s:17:"attrib_color_name";}s:18:"_1447682983606_606";a:4:{s:5:"order";s:3:"220";s:6:"column";s:4:"test";s:9:"attribute";s:37:"xunseria_directive_variant_attributes";s:5:"param";a:5:{i:0;s:11:"attrib_test";i:1;s:19:"attrib_test_general";i:2;s:14:"attrib_test_cm";i:3;s:20:"attrib_swim_cup_test";i:4;s:21:"attrib_width_cup_test";}}s:18:"_1447682998765_765";a:3:{s:5:"order";s:3:"240";s:6:"column";s:13:"item_group_id";s:9:"attribute";s:20:"attrib_articlenumber";}s:18:"_1447683031972_972";a:4:{s:5:"order";s:3:"250";s:6:"column";s:6:"gender";s:9:"attribute";s:37:"xunseria_directive_variant_attributes";s:5:"param";a:1:{i:0;s:13:"attrib_gender";}}s:18:"_1447683060883_883";a:4:{s:5:"order";s:3:"260";s:6:"column";s:9:"age_group";s:9:"attribute";s:31:"xunseria_directive_static_value";s:5:"param";s:5:"adult";}s:18:"_1447683070403_403";a:3:{s:5:"order";s:3:"270";s:6:"column";s:8:"material";s:9:"attribute";s:37:"xunseria_directive_variant_attributes";}s:18:"_1447683088418_418";a:3:{s:5:"order";s:3:"280";s:6:"column";s:4:"gtin";s:9:"attribute";s:3:"sku";}}`

var toSerialize = phpserialize.PhpArray{
	"_1447683031972_972": phpserialize.PhpArray{"order": "250", "column": "gender", "attribute": "xunseria_directive_variant_attributes",
		"param": phpserialize.PhpArray{0: "attrib_gender"}},
	"_1447683070403_403": phpserialize.PhpArray{"order": "270", "column": "material", "attribute": "xunseria_directive_variant_attributes"},
	2:                    phpserialize.PhpArray{"order": "20", "column": "title", "attribute": "xunseria_directive_concatenate", "param": "attrib {{attrib_name_map}} {{attrib_serie}} {{name}} {{attrib_articlenumber}}"},
	3:                    phpserialize.PhpArray{"order": "30", "column": "description", "attribute": "description"},
	7:                    phpserialize.PhpArray{"order": "70", "column": "price", "attribute": "xunseria_directive_price", "param": "1"},
	9:                    phpserialize.PhpArray{"order": "90", "column": "sale_price_effective_date", "attribute": "xunseria_directive_sale_price_effective_date"},
	10:                   phpserialize.PhpArray{"order": "110", "column": "availability", "attribute": "xunseria_directive_availability"},
	13:                   phpserialize.PhpArray{"order": "150", "column": "mpn", "attribute": "attrib_articlenumber"},
	"_1447683088418_418": phpserialize.PhpArray{"order": "280", "column": "gtin", "attribute": "sku"},
	"_1447682983606_606": phpserialize.PhpArray{"order": "220", "column": "test", "attribute": "xunseria_directive_variant_attributes",
		"param": phpserialize.PhpArray{3: "attrib_swim_cup_test", 4: "attrib_width_cup_test", 0: "attrib_test", 1: "attrib_test_general", 2: "attrib_test_cm"}},
	"_1447682998765_765": phpserialize.PhpArray{"order": "240", "column": "item_group_id", "attribute": "attrib_articlenumber"},
	0:                    phpserialize.PhpArray{"order": "10", "column": "id", "attribute": "xunseria_directive_id", "param": "0"},
	4:                    phpserialize.PhpArray{"order": "40", "column": "link", "attribute": "xunseria_directive_url", "param": "?utm_source=google_shopping"},
	5:                    phpserialize.PhpArray{"order": "50", "column": "image_link", "attribute": "xunseria_directive_image_link", "param": "image"},
	12:                   phpserialize.PhpArray{"param": "attrib", "order": "140", "column": "brand", "attribute": "xunseria_directive_static_value"},
	14:                   phpserialize.PhpArray{"column": "condition", "attribute": "xunseria_directive_static_value", "param": "new", "order": "160"},
	"_1447682911167_167": phpserialize.PhpArray{"order": "210", "column": "color", "attribute": "attrib_color_name"},
	1:                    phpserialize.PhpArray{"order": "10", "column": "item_group_id", "attribute": "xunseria_directive_item_group_id"},
	8:                    phpserialize.PhpArray{"attribute": "xunseria_directive_sale_price", "param": "1", "order": "80", "column": "sale_price"},
	15:                   phpserialize.PhpArray{"order": "170", "column": "product_type", "attribute": "xunseria_directive_product_type_by_category"},
	"_1447683060883_883": phpserialize.PhpArray{"order": "260", "column": "age_group", "attribute": "xunseria_directive_static_value", "param": "adult"},
	6:                    phpserialize.PhpArray{"param": "image", "order": "60", "column": "additional_image_link", "attribute": "xunseria_directive_additional_image_link"},
	16:                   phpserialize.PhpArray{"column": "google_product_category", "attribute": "xunseria_directive_google_category_by_category", "order": "180"},
	17:                   phpserialize.PhpArray{"attribute": "xunseria_directive_identifier_exists", "order": "190", "column": "identifier_exists"},
	18:                   phpserialize.PhpArray{"order": "200", "column": "is_bundle", "attribute": "xunseria_directive_is_bundle"},
}
