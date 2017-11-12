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

package i18n_test

import (
	"testing"

	"github.com/corestoreio/pkg/i18n"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/stretchr/testify/assert"
)

func TestLocaleAvailable(t *testing.T) {
	want := slices.String{"aa_ET", "ab_GE", "ae_IR", "af_ZA", "ak_GH", "am_ET", "ar_EG", "as_IN", "av_RU", "ay_BO", "az_AZ", "ba_RU", "be_BY", "bg_BG", "bi_VU", "bm_ML", "bn_BD", "bo_CN", "br_FR", "bs_BA", "ca_ES", "ce_RU", "ch_GU", "co_FR", "cr_CA", "cs_CZ", "cu_RU", "cv_RU", "cy_GB", "da_DK", "de_DE", "dv_MV", "dz_BT", "ee_GH", "el_GR", "en_US", "es_ES", "et_EE", "eu_ES", "fa_IR", "ff_SN", "fi_FI", "fj_FJ", "fo_FO", "fr_FR", "fy_NL", "ga_IE", "gd_GB", "gl_ES", "gn_PY", "gu_IN", "gv_IM", "ha_NG", "he_IL", "hi_IN", "ho_PG", "hr_HR", "ht_HT", "hu_HU", "hy_AM", "hz_NA", "ia_FR", "id_ID", "ig_NG", "ii_CN", "ik_US", "is_IS", "it_IT", "iu_CA", "ja_JP", "jv_ID", "ka_GE", "kg_CD", "ki_KE", "kj_NA", "kk_KZ", "kl_GL", "km_KH", "kn_IN", "ko_KR", "ks_IN", "ku_TR", "kv_RU", "kw_GB", "ky_KG", "la_VA", "lb_LU", "lg_UG", "li_NL", "ln_CD", "lo_LA", "lt_LT", "lu_CD", "lv_LV", "mg_MG", "mh_MH", "mi_NZ", "mk_MK", "ml_IN", "mn_MN", "mr_IN", "ms_MY", "mt_MT", "my_MM", "na_NR", "nd_ZW", "ne_NP", "ng_NA", "nl_NL", "nn_NO", "no_NO", "nr_ZA", "nv_US", "ny_MW", "oc_FR", "om_ET", "or_IN", "os_GE", "pa_IN", "pl_PL", "ps_AF", "pt_BR", "qu_PE", "rm_CH", "rn_BI", "ro_RO", "ru_RU", "rw_RW", "sa_IN", "sc_IT", "sd_PK", "se_NO", "sg_CF", "si_LK", "sk_SK", "sl_SI", "sm_WS", "sn_ZW", "so_SO", "sq_AL", "sr_RS", "ss_ZA", "st_ZA", "su_ID", "sv_SE", "sw_TZ", "ta_IN", "te_IN", "tg_TJ", "th_TH", "ti_ET", "tk_TM", "tn_ZA", "to_TO", "tr_TR", "ts_ZA", "tt_RU", "ty_PF", "ug_CN", "uk_UA", "ur_PK", "uz_UZ", "ve_ZA", "vi_VN", "wa_BE", "wo_SN", "xh_ZA", "yo_NG", "za_CN", "zh_CN", "zu_ZA", "ace_ID", "ach_UG", "ada_GH", "ady_RU", "aeb_TN", "agq_CM", "akk_IQ", "alt_RU", "arc_IR", "arn_CL", "aro_BO", "arq_DZ", "ary_MA", "arz_EG", "asa_TZ", "ase_US", "ast_ES", "awa_IN", "bal_PK", "ban_ID", "bar_AT", "bas_CM", "bax_CM", "bbc_ID", "bbj_CM", "bej_SD", "bem_ZM", "bew_ID", "bez_TZ", "bfd_CM", "bfq_IN", "bgn_PK", "bho_IN", "bik_PH", "bin_NG", "bjn_ID", "bkm_CM", "bpy_IN", "bqi_IR", "bra_IN", "brh_PK", "brx_IN", "bss_CM", "bua_RU", "bug_ID", "bum_CM", "byn_ER", "byv_CM", "cch_NG", "ceb_PH", "cgg_UG", "chk_FM", "chm_RU", "cho_US", "chp_CA", "chr_US", "ckb_IQ", "cop_EG", "cps_PH", "csb_PL", "dak_US", "dar_RU", "dav_KE", "den_CA", "dgr_CA", "dje_NE", "doi_IN", "dsb_DE", "dtp_MY", "dua_CM", "dyo_SN", "dyu_BF", "ebu_KE", "efi_NG", "egl_IT", "egy_EG", "esu_US", "ewo_CM", "ext_ES", "fan_GQ", "fil_PH", "fit_SE", "fon_BJ", "frc_US", "frp_FR", "frr_DE", "frs_DE", "fur_IT", "gaa_GH", "gag_MD", "gan_CN", "gay_ID", "gbz_IR", "gez_ET", "gil_KI", "glk_IR", "gom_IN", "gon_IN", "gor_ID", "got_UA", "grc_CY", "gsw_CH", "guc_CO", "gur_GH", "guz_KE", "gwi_CA", "hak_CN", "haw_US", "hif_FJ", "hil_PH", "hsb_DE", "hsn_CN", "iba_MY", "ibb_NG", "ilo_PH", "inh_RU", "izh_RU", "jam_JM", "jgo_CM", "jmc_TZ", "jut_DK", "kaa_UZ", "kab_DZ", "kac_MM", "kaj_NG", "kam_KE", "kbd_RU", "kcg_NG", "kde_TZ", "kea_CV", "ken_CM", "kfo_CI", "kgp_BR", "kha_IN", "khq_ML", "khw_PK", "kiu_TR", "kkj_CM", "kln_KE", "kmb_AO", "koi_RU", "kok_IN", "kos_FM", "kpe_LR", "krc_RU", "kri_SL", "krj_PH", "krl_RU", "kru_IN", "ksb_TZ", "ksf_CM", "ksh_DE", "kum_RU", "lad_IL", "lag_TZ", "lah_PK", "lez_RU", "lij_IT", "lkt_US", "lmo_IT", "lol_CD", "loz_ZM", "lrc_IR", "ltg_LV", "lua_CD", "luo_KE", "luy_KE", "lzh_CN", "lzz_TR", "mad_ID", "maf_CM", "mag_IN", "mai_IN", "mak_ID", "man_GM", "mas_KE", "mdf_RU", "mdr_ID", "men_SL", "mer_KE", "mfe_MU", "mgh_MZ", "mgo_CM", "min_ID", "mni_IN", "moh_CA", "mos_BF", "mrj_RU", "mua_CM", "mus_US", "mwr_IN", "mwv_ID", "myv_RU", "mzn_IR", "nan_CN", "nap_IT", "naq_NA", "nds_DE", "new_NP", "niu_NU", "njo_IN", "nmg_CM", "nnh_CM", "non_SE", "nqo_GN", "nso_ZA", "nus_SS", "nym_TZ", "nyn_UG", "nzi_GH", "pag_PH", "pal_IR", "pam_PH", "pap_AW", "pau_PW", "pcd_FR", "pdc_US", "pdt_CA", "peo_IR", "pfl_DE", "phn_LB", "pms_IT", "pnt_GR", "pon_FM", "quc_GT", "qug_EC", "raj_IN", "rgn_IT", "rif_MA", "rof_TZ", "rtm_FJ", "rue_UA", "rug_SB", "rwk_TZ", "sah_RU", "saq_KE", "sas_ID", "sat_IN", "saz_IN", "sbp_TZ", "scn_IT", "sco_GB", "sdc_IT", "sdh_IR", "seh_MZ", "sei_MX", "ses_ML", "sga_IE", "sgs_LT", "shi_MA", "shn_MM", "sid_ET", "sli_PL", "sly_ID", "sma_SE", "smj_SE", "smn_FI", "sms_FI", "snk_ML", "srn_SR", "srr_SN", "ssy_ER", "stq_DE", "suk_TZ", "sus_GN", "swb_YT", "swc_CD", "syr_IQ", "szl_PL", "tcy_IN", "tem_SL", "teo_UG", "tet_TL", "tig_ER", "tiv_NG", "tkl_TK", "tkr_AZ", "tly_AZ", "tmh_NE", "tog_MW", "tpi_PG", "tru_TR", "trv_TW", "tsd_GR", "ttt_AZ", "tum_MW", "tvl_TV", "twq_NE", "tyv_RU", "tzm_MA", "udm_RU", "uga_SY", "umb_AO", "vai_LR", "vec_IT", "vep_RU", "vls_BE", "vmf_DE", "vot_RU", "vro_EE", "vun_TZ", "wae_CH", "wal_ET", "war_PH", "wbp_AU", "wuu_CN", "xmf_GE", "xog_UG", "yao_MZ", "yap_FM", "yav_CM", "ybb_CM", "yrl_BR", "zea_NL", "zgh_MA", "zza_TR", "az_IR", "de_AT", "de_CH", "en_AU", "en_CA", "en_GB", "en_US", "es_ES", "es_MX", "fa_AF", "fr_CA", "fr_CH", "nds_NL", "nl_BE", "pt_BR", "pt_PT", "ro_MD", "sr_RS", "zh_CN", "zh_TW"}
	assert.EqualValues(t, want, i18n.LocaleAvailable)
}

func TestLocaleSupported(t *testing.T) {
	want := slices.String{"en_US", "de_DE", "fr_FR", "it_IT", "es_ES", "ja_JP", "uk_UA"}
	assert.EqualValues(t, want, i18n.LocaleSupported)
}
