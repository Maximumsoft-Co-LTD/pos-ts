# pos-ts
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o printslip.exe _cmd/main.go

			// p.SetSmooth(1)
			// p.SetFontSize(2, 3)
			// p.SetFont("A")
			// p.Write("test ")
			// p.SetFont("B")
			// p.Write("test2 ")
			// p.SetFont("C")
			// p.Write("test3 ")
			// p.Formfeed()

			// p.SetFont("B")
			// p.SetFontSize(1, 1)

			// p.SetEmphasize(1)
			// p.Write("halle")
			// p.Formfeed()

			// p.SetUnderline(1)
			// p.SetFontSize(4, 4)
			// p.Write("halle")

			// p.SetReverse(1)
			// p.SetFontSize(2, 4)
			// p.Write("halle")
			// p.Formfeed()

			// p.SetFont("C")
			// p.SetFontSize(8, 8)
			// p.Write("halle")
			// p.FormfeedN(5)

			// Set font, style, and size