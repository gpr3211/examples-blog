package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	"strings"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type DebugUI struct {
	UI                  debugui.DebugUI
	inputCapturingState debugui.InputCapturingState
	Perm                *DebugLog
	DebugMode           bool
	FastGame            bool
	Temp                *TempLog
	HiRes               bool
	Lang                bool
}

// TempLog for in game debugging
type TempLog struct {
	logUpdated   bool
	logSubmitBuf string
	logBuf       string
}

// DebugLog perma log saved to file
type DebugLog struct {
	logUpdated   bool
	logSubmitBuf string
	logBuf       string
	logData      []string
	logfile      *os.File
}

func NewDebugUI() *DebugUI {

	g := &DebugUI{
		Temp: &TempLog{},
		Perm: &DebugLog{},
		Lang: false,
	}

	writeInit := func() {
		f, err := os.Open("debug_log.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		lines := getLinesChannel(f)
		for v := range lines {
			if len(g.Perm.logBuf) > 0 {
				g.Perm.logBuf += "\n"
			}
			g.Perm.logData = append(g.Perm.logData, (fmt.Sprintf("%s", v)))
			g.Perm.logBuf += (fmt.Sprintf("%s", v))
			g.Perm.logUpdated = true
		}
	}
	writeInit()

	return g
}

func (d *DebugUI) Update() error {
	inputcap, _ := d.UI.Update(func(ctx *debugui.Context) error {
		d.logWindow(ctx)
		d.templogWindow(ctx)
		return nil
	})
	d.inputCapturingState = inputcap
	return nil
}

func (g *DebugUI) logWindow(ctx *debugui.Context) {
	ctx.Window("Perma Log", image.Rect(50, 50, 450, 150), func(layout debugui.ContainerLayout) {
		ctx.SetGridLayout([]int{-1}, []int{-1, 0})
		ctx.Panel(func(layout debugui.ContainerLayout) {
			ctx.SetGridLayout([]int{-1}, []int{-1})
			ctx.Text(g.Perm.logBuf)
			if g.Perm.logUpdated {
				ctx.SetScroll(image.Pt(layout.ScrollOffset.X, layout.ContentSize.Y))
				g.Perm.logUpdated = false
			}
		})
		ctx.GridCell(func(bounds image.Rectangle) {
			submit := func() {
				if g.Perm.logSubmitBuf == "" {
					return
				}
				g.writePermLog(g.Perm.logSubmitBuf)
				g.Perm.logSubmitBuf = ""
			}
			ctx.SetGridLayout([]int{-3, -1}, nil)
			ctx.TextField(&g.Perm.logSubmitBuf).On(func() {
				if ebiten.IsKeyPressed(ebiten.KeyEnter) {
					submit()
					ctx.SetTextFieldValue(g.Perm.logSubmitBuf)
				}
			})
			ctx.Button("Submit").On(func() {
				submit()
			})
		})
	})
}

// writeLog writes a permanent debug log to debug_log.txt and shows

func (g *DebugUI) writePermLog(text string) {
	if len(g.Perm.logBuf) > 0 {
		g.Perm.logBuf += "\n"
	}
	split := strings.Split(text, " ")
	if len(split) > 0 {
		if split[0] == "remove" {

		}
	}
	g.Perm.logBuf += text
	g.Perm.logUpdated = true
	f := "debug_log.txt"
	f4, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f4.Close()
	_, err = f4.Write([]byte(text + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// getLinesChannel returns a channel that shoots the msg line by line.
func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)
	go func() {
		r := bufio.NewReader(f)
		line := []byte{}
		line_index := 0
		for {
			buff := make([]byte, 8)
			chunk, _ := r.Read(buff)
			buff = buff[:chunk]
			for _, v := range buff {
				if rune(v) != '\n' {
					line = append(line, v)
				} else {
					linesChan <- fmt.Sprintf("%d:", line_index) + string(line)
					line_index++
					line = []byte{}
				}
			}
			if chunk == 0 {
				if len(line) > 0 {
					linesChan <- string(line)
				}
				close(linesChan)
				return
			}
		}
	}()
	return linesChan
}

// TEMP LOG
func (g *DebugUI) templogWindow(ctx *debugui.Context) {
	ctx.Window("Log", image.Rect(350, 300, 650, 500), func(layout debugui.ContainerLayout) {
		ctx.Header("Game Config", true, func() {
			ctx.Checkbox(&g.HiRes, "Hi-Res").On(func() {
				if g.HiRes {
					ctx.SetScale(2)
				} else {
					ctx.SetScale(1)
				}
			})
			ctx.Checkbox(&g.Lang, "Lang").On(func() {
			})

		})
		ctx.SetGridLayout([]int{-1}, []int{-1, 0})
		ctx.Panel(func(layout debugui.ContainerLayout) {
			ctx.SetGridLayout([]int{-1}, []int{-1})
			ctx.Text(g.Temp.logBuf)
			if g.Temp.logUpdated {
				ctx.SetScroll(image.Pt(layout.ScrollOffset.X, layout.ContentSize.Y))
				g.Temp.logUpdated = false
			}
		})
		ctx.GridCell(func(bounds image.Rectangle) {
			submit := func() {
				if g.Temp.logSubmitBuf == "" {
					return
				}
				g.WriteLog(g.Temp.logSubmitBuf)
				g.Temp.logSubmitBuf = ""
			}

			ctx.SetGridLayout([]int{-3, -1}, nil)
			ctx.TextField(&g.Temp.logSubmitBuf).On(func() {
				if ebiten.IsKeyPressed(ebiten.KeyEnter) {
					submit()
					ctx.SetTextFieldValue(g.Temp.logSubmitBuf)
				}
			})
			ctx.Button("Submit").On(func() {
				submit()
			})
		})
	})
}

func (g *DebugUI) WriteLog(text string) {
	if len(g.Temp.logBuf) > 0 {
		g.Temp.logBuf += "\n"
	}
	g.Temp.logBuf += text
	g.Temp.logUpdated = true
}
