package main

import (
	sfCsv "github.com/t8nlab/surface/csv"
	sfImage "github.com/t8nlab/surface/image"
	sfJson "github.com/t8nlab/surface/json"
	sfSmtp "github.com/t8nlab/surface/smtp"
	sfClean "github.com/t8nlab/surface/clean"
	sfExtract "github.com/t8nlab/surface/extract"
	ext "github.com/t8nlab/surface/utils"
)

func init() {
	ext.Register("csv_open", sfCsv.CsvOpen)
	ext.Register("csv_next", sfCsv.CsvNext)
	ext.Register("csv_read_all", sfCsv.CsvReadAll)
	ext.Register("csv_close", sfCsv.CsvClose)
	ext.Register("csv_create", sfCsv.CsvCreate)
	ext.Register("csv_write", sfCsv.CsvWrite)

	ext.Register("smtp_send", sfSmtp.SmtpSend)
	ext.Register("smtp_bulk_send", sfSmtp.SmtpBulkSend)
	ext.Register("smtp_render", sfSmtp.SmtpRender)
	ext.Register("smtp_render_file", sfSmtp.SmtpRenderFile)

	ext.Register("image_resize", sfImage.ImageResize)
	ext.Register("image_crop", sfImage.ImageCrop)
	ext.Register("image_process", sfImage.ImageProcess)
	ext.Register("image_batch", sfImage.ImageBatch)

	ext.Register("json_open", sfJson.JsonOpen)
	ext.Register("json_next", sfJson.JsonNext)
	ext.Register("json_read_all", sfJson.JsonReadAll)
	ext.Register("json_close", sfJson.JsonClose)
	ext.Register("json_create", sfJson.JsonCreate)
	ext.Register("json_write", sfJson.JsonWrite)
	ext.Register("json_stringify", sfJson.JsonStringifyFast)
	ext.Register("json_to_csv", sfJson.JsonToCsv)

	ext.Register("clean_validate_emails", sfClean.ValidateEmails)
	ext.Register("clean_normalize_phones", sfClean.NormalizePhones)
	ext.Register("clean_remove_duplicates", sfClean.RemoveDuplicates)
	ext.Register("clean_process", sfClean.Process)

	ext.Register("extract_html", sfExtract.ExtractHTML)
	ext.Register("extract_links", sfExtract.ExtractLinks)
	ext.Register("extract_meta", sfExtract.ExtractMeta)
}


func main() {}