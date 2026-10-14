"""Genera el PDF del informe medido; no modifica sus datos."""
from pathlib import Path
import re
from xml.sax.saxutils import escape

from reportlab.lib import colors
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer, PageBreak, Table, TableStyle
from pypdf import PdfReader

root = Path(__file__).resolve().parents[1]
source = root / "docs" / "INFORME_OPERACION.md"
output = root / "docs" / "INFORME_OPERACION.pdf"
styles = getSampleStyleSheet()
styles.add(ParagraphStyle(name="P10Body", fontName="Helvetica", fontSize=10, leading=14, spaceAfter=8))
styles["Heading1"].fontSize = 18
styles["Heading2"].fontSize = 12
styles["Heading2"].spaceBefore = 8
styles["Heading2"].spaceAfter = 6
story = []


def fmt(text):
    text = escape(text)
    return re.sub(r"\*\*(.*?)\*\*", r"<b>\1</b>", text)


for block in source.read_text(encoding="utf-8").split("\n\n"):
    if not block.strip():
        continue
    if block.startswith("# "):
        story.append(Paragraph(fmt(block[2:]), styles["Heading1"]))
    elif block.startswith("## "):
        if block.startswith("## 4."):
            story.append(PageBreak())
        story.append(Paragraph(fmt(block[3:]), styles["Heading2"]))
    elif block.startswith("|"):
        rows = []
        for line in block.splitlines():
            if "---" in line:
                continue
            rows.append([Paragraph(fmt(cell.strip()), styles["P10Body"]) for cell in line.strip("|").split("|")])
        table = Table(rows, colWidths=[330, 165], hAlign="LEFT")
        table.setStyle(TableStyle([
            ("BACKGROUND", (0, 0), (-1, 0), colors.HexColor("#e9edf2")),
            ("GRID", (0, 0), (-1, -1), 0.5, colors.HexColor("#c4ccd5")),
            ("VALIGN", (0, 0), (-1, -1), "TOP"),
            ("BOTTOMPADDING", (0, 0), (-1, -1), 4),
        ]))
        story += [table, Spacer(1, 8)]
    else:
        story.append(Paragraph(fmt(" ".join(block.splitlines())), styles["P10Body"]))


def footer(canvas, doc):
    canvas.setFont("Helvetica", 9)
    canvas.drawString(42, 24, "Practica 10 - 202300644 - datos medidos el 14/10/2026")
    canvas.drawRightString(A4[0]-42, 24, str(doc.page))


SimpleDocTemplate(str(output), pagesize=A4, leftMargin=42, rightMargin=42,
                  topMargin=38, bottomMargin=42).build(story, onFirstPage=footer, onLaterPages=footer)
pages = len(PdfReader(output).pages)
if pages > 2:
    raise RuntimeError(f"El informe excede el limite: {pages} paginas")
print(f"Informe generado: {pages} paginas")
