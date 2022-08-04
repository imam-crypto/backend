package helper

import "github.com/jung-kurt/gofpdf"

func SetToPDF() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetLeftMargin(30)
	pdf.SetFont("Times", "B", 14)
	pdf.Cell(40, 10, "Data User")
	pdf.Ln(12)

	return pdf
}
func Header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 13)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		pdf.CellFormat(49, 7, str, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}
func Table(pdf *gofpdf.Fpdf, data [][]string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 12)
	pdf.SetFillColor(255, 255, 255)
	//pdf.SetLeftMargin(20)
	align := []string{"L", "C", "L", "R", "R", "R"}

	for _, line := range data {
		for i, str := range line {
			pdf.CellFormat(49, 7, str, "1", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}
func SaveFile(pdf *gofpdf.Fpdf) error {

	return pdf.OutputFileAndClose("./file/DataUser.pdf")
}
