package render

import (
	"fmt"
	"log/slog"
	"strings"

	extast "github.com/yuin/goldmark/extension/ast"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.Prioritized(newTelegramRenderer(), 500),
				),
			)),
		goldmark.WithExtensions(extension.Table, extension.Strikethrough),
	)
}

func AdjustMdToTelegramFormat(input string) (string, error) {
	var buf strings.Builder
	if err := md.Convert([]byte(input), &buf); err != nil {
		slog.Error("unable to convert input to tg format", slog.Any("error", err), slog.String("input", input))
		return "", fmt.Errorf("unable to adjust input: %w", err)
	}

	result := buf.String()
	result = strings.TrimSpace(result)
	return result, nil
}

type TelegramRenderer struct{}

func newTelegramRenderer(opts ...goldmark.Option) renderer.NodeRenderer {
	return &TelegramRenderer{}
}

func escapeTelegramMarkdown(text string) string {
	var b strings.Builder
	b.Grow(len(text) * 2)

	for _, r := range text {
		switch r {
		case '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!':
			b.WriteByte('\\')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (r *TelegramRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(extast.KindStrikethrough, r.renderStrikethrough)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)

	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindText, r.renderText)

	reg.Register(extast.KindTable, r.renderTable)
	reg.Register(extast.KindTableHeader, r.renderTableHeader)
	reg.Register(extast.KindTableRow, r.renderTableRow)
	reg.Register(extast.KindTableCell, r.renderTableCell)

	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
}

func (r *TelegramRenderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n---\n")
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n> ")
	} else {
		w.WriteString("\n")
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	w.WriteString("~")
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	w.WriteString("`")
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n**tables not supported**\n")
	}
	return ast.WalkSkipChildren, nil
}

func (r *TelegramRenderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	emphasis := node.(*ast.Emphasis)

	parentIsEmphasis := false
	if parent := node.Parent(); parent != nil {
		if parentEmph, ok := parent.(*ast.Emphasis); ok {
			parentIsEmphasis = true
			if entering {
				if parentEmph.Level == 1 && emphasis.Level == 2 {
					w.WriteString("_")
				} else if emphasis.Level == 2 {
					w.WriteString("*")
				} else {
					w.WriteString("_")
				}
			} else {
				if parentEmph.Level == 1 && emphasis.Level == 2 {
					w.WriteString("_")
				} else if emphasis.Level == 2 {
					w.WriteString("*")
				} else {
					w.WriteString("_")
				}
			}
			return ast.WalkContinue, nil
		}
	}

	hasChildEmphasis := false
	if emphasis.ChildCount() > 0 {
		for child := emphasis.FirstChild(); child != nil; child = child.NextSibling() {
			if _, ok := child.(*ast.Emphasis); ok {
				hasChildEmphasis = true
				break
			}
		}
	}

	if !parentIsEmphasis {
		if entering {
			if hasChildEmphasis {
				w.WriteString("*")
			} else if emphasis.Level == 2 {
				w.WriteString("*")
			} else {
				w.WriteString("_")
			}
		} else {
			if hasChildEmphasis {
				w.WriteString("*")
			} else if emphasis.Level == 2 {
				w.WriteString("*")
			} else {
				w.WriteString("_")
			}
		}
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	if entering {
		lang := string(n.Language([]byte(source)))
		if lang != "" {
			fmt.Fprintf(w, "```%s\n", lang)
		} else {
			w.WriteString("```\n")
		}

		var code strings.Builder
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			code.Write(line.Value([]byte(source)))
		}

		w.WriteString(code.String())
	} else {
		w.WriteString("```")
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.String)
	if entering {
		w.WriteString(escapeTelegramMarkdown(string(n.Value)))
	}
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	destination := string(n.Destination)
	title := string(n.Title)
	if entering {
		w.WriteString("[")
	} else {
		w.WriteString("](")
		w.WriteString(destination)
		if len(title) > 0 {
			w.WriteString(" \"")
			w.WriteString(title)
			w.WriteString("\"")
		}
		w.WriteString(")")
	}
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderBlockSeparatorsAround(renderer func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error)) func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		_, err := r.renderBlockSeparator(w, source, node, entering)
		if err != nil {
			return ast.WalkStop, err
		}
		_, err = renderer(w, source, node, entering)
		if err != nil {
			return ast.WalkStop, err
		}
		_, err = r.renderBlockSeparator(w, source, node, entering)
		if err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkContinue, nil
	}
}

func (r *TelegramRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n*")
	} else {
		w.WriteString("*\n")
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Text)
	if entering {
		text := n.Value(source)

		w.WriteString(escapeTelegramMarkdown(string(text)))
		if n.SoftLineBreak() {
			w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderWalkContinue(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderBlockSeparator(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if node.PreviousSibling() != nil && node.HasBlankPreviousLines() {
			w.WriteString("\n")
		}
	} else {
		w.Flush()
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		n := node.(*ast.Paragraph)
		if n.Parent().Kind() != ast.KindListItem && n.Parent().Kind() != ast.KindBlockquote {
			if n.NextSibling() != nil {
				w.WriteString("\n")
			}
		}
	}

	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *TelegramRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.ListItem)
		parent := n.Parent().(*ast.List)

		if parent.IsOrdered() {
			itemNum := 1
			prev := n.PreviousSibling()
			for prev != nil {
				itemNum++
				prev = prev.PreviousSibling()
			}
			fmt.Fprintf(w, "%d\\. ", itemNum)
		} else {
			w.WriteString("• ")
		}
	} else {
		w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}
