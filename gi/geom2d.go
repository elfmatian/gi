// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"fmt"
	"image"
	"log"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/chewxy/math32"
	"github.com/goki/ki"
	"github.com/goki/ki/kit"

	// "github.com/rcoreilly/rasterx"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

// SVG default coordinates are such that 0,0 is upper-left!

/*
This is heavily modified from: https://github.com/fogleman/gg

Copyright (C) 2016 Michael Fogleman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// note: golang.org/x/image/math/f64 defines Vec2 as [2]float64
// elabored then by https://godoc.org/github.com/go-gl/mathgl/mgl64
// it is instead very convenient and clear to use .X .Y fields for 2D math
// original gg package used Point2D but Vec2D is more general, e.g., for sizes etc
// in go much better to use fewer types so only using Vec2D

// could break this out as separate package, but no advantage in package-based
// naming

////////////////////////////////////////////////////////////////////////////////////////
//  Min / Max for other types..

// math provides Max/Min for 64bit -- these are for specific subtypes

func Max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// SetMax32 sets arg a to Max(a,b)
func SetMax32(a *float32, b float32) {
	if *a < b {
		*a = b
	}
}

func Min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// SetMin32 sets arg a to Min(a,b)
func SetMin32(a *float32, b float32) {
	if *a > b {
		*a = b
	}
}

// MinPos returns the minimum of the two values, excluding any that are 0
func MinPos(a, b float64) float64 {
	if a > 0.0 && b > 0.0 {
		return math.Min(a, b)
	} else if a > 0.0 {
		return a
	} else if b > 0.0 {
		return b
	}
	return a
}

// MinPos32 returns the minimum of the two values, excluding any that are 0
func MinPos32(a, b float32) float32 {
	if a > 0.0 && b > 0.0 {
		return Min32(a, b)
	} else if a > 0.0 {
		return a
	} else if b > 0.0 {
		return b
	}
	return a
}

// InRange returns the value constrained to the min / max range
func InRange(val, min, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// InRange32 returns the value constrained to the min / max range
func InRange32(val, min, max float32) float32 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// InRangeInt returns the value constrained to the min / max range
func InRangeInt(val, min, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// Truncate a floating point number to given level of precision -- slow.. uses string conversion
func Truncate(val float64, prec int) float64 {
	frep := strconv.FormatFloat(val, 'g', prec, 64)
	val, _ = strconv.ParseFloat(frep, 64)
	return val
}

// Truncate a floating point number to given level of precision -- slow.. uses string conversion
func Truncate32(val float32, prec int) float32 {
	frep := strconv.FormatFloat(float64(val), 'g', prec, 32)
	tval, _ := strconv.ParseFloat(frep, 32)
	return float32(tval)
}

// FloatMod ensures that a floating point value is an even multiple of a given value
func FloatMod(val, mod float64) float64 {
	return float64(int(math.Round(val/mod))) * mod
}

// FloatMod ensures that a floating point value is an even multiple of a given value
func FloatMod32(val, mod float32) float32 {
	return float32(int(math.Round(float64(val/mod)))) * mod
}

// dimensions
type Dims2D int32

const (
	X Dims2D = iota
	Y
	Dims2DN
)

var KiT_Dims2D = kit.Enums.AddEnumAltLower(Dims2DN, false, nil, "")

func (ev Dims2D) MarshalJSON() ([]byte, error)  { return kit.EnumMarshalJSON(ev) }
func (ev *Dims2D) UnmarshalJSON(b []byte) error { return kit.EnumUnmarshalJSON(ev, b) }

//go:generate stringer -type=Dims2D

// 2D vector -- a point or size in 2D
type Vec2D struct {
	X, Y float32
}

var Vec2DZero = Vec2D{0, 0}

func NewVec2D(x, y float32) Vec2D {
	return Vec2D{x, y}
}

func NewVec2DFmPoint(pt image.Point) Vec2D {
	v := Vec2D{}
	v.SetPoint(pt)
	return v
}

func NewVec2DFmFixed(pt fixed.Point26_6) Vec2D {
	v := Vec2D{}
	v.SetFixed(pt)
	return v
}

// return value along given dimension
func (a Vec2D) Dim(d Dims2D) float32 {
	switch d {
	case X:
		return a.X
	default:
		return a.Y
	}
}

// get the other dimension
func OtherDim(d Dims2D) Dims2D {
	switch d {
	case X:
		return Y
	default:
		return X
	}
}

// set the value along a given dimension
func (a *Vec2D) SetDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X = val
	case Y:
		a.Y = val
	}
}

// set values
func (a *Vec2D) Set(x, y float32) {
	a.X = x
	a.Y = y
}

// set both dims to same value
func (a *Vec2D) SetVal(val float32) {
	a.X = val
	a.Y = val
}

func (a Vec2D) IsZero() bool {
	return a.X == 0.0 && a.Y == 0.0
}

func (a Vec2D) Fixed() fixed.Point26_6 {
	return Float32ToFixedPoint(a.X, a.Y)
}

func (a Vec2D) Add(b Vec2D) Vec2D {
	return Vec2D{a.X + b.X, a.Y + b.Y}
}

func (a Vec2D) AddVal(val float32) Vec2D {
	return Vec2D{a.X + val, a.Y + val}
}

func (a *Vec2D) SetAdd(b Vec2D) {
	a.X += b.X
	a.Y += b.Y
}

func (a *Vec2D) SetAddVal(val float32) {
	a.X += val
	a.Y += val
}

func (a *Vec2D) SetAddDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X += val
	case Y:
		a.Y += val
	}
}

func (a Vec2D) Sub(b Vec2D) Vec2D {
	return Vec2D{a.X - b.X, a.Y - b.Y}
}

func (a *Vec2D) SetSub(b Vec2D) {
	a.X -= b.X
	a.Y -= b.Y
}

func (a *Vec2D) SetSubVal(val float32) {
	a.X -= val
	a.Y -= val
}

func (a *Vec2D) SetSubDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X -= val
	case Y:
		a.Y -= val
	}
}

func (a Vec2D) SubVal(val float32) Vec2D {
	return Vec2D{a.X - val, a.Y - val}
}

func (a Vec2D) Mul(b Vec2D) Vec2D {
	return Vec2D{a.X * b.X, a.Y * b.Y}
}

func (a *Vec2D) SetMul(b Vec2D) {
	a.X *= b.X
	a.Y *= b.Y
}

func (a Vec2D) MulVal(val float32) Vec2D {
	return Vec2D{a.X * val, a.Y * val}
}

func (a *Vec2D) SetMulVal(val float32) {
	a.X *= val
	a.Y *= val
}

func (a *Vec2D) SetMulDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X *= val
	case Y:
		a.Y *= val
	}
}

func (a Vec2D) Div(b Vec2D) Vec2D {
	return Vec2D{a.X / b.X, a.Y / b.Y}
}

func (a *Vec2D) SetDiv(b Vec2D) {
	a.X /= b.X
	a.Y /= b.Y
}

func (a *Vec2D) SetDivlVal(val float32) {
	a.X /= val
	a.Y /= val
}

func (a *Vec2D) SetDivDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X /= val
	case Y:
		a.Y /= val
	}
}

func (a Vec2D) DivVal(val float32) Vec2D {
	return Vec2D{a.X / val, a.Y / val}
}

func (a Vec2D) Max(b Vec2D) Vec2D {
	return Vec2D{Max32(a.X, b.X), Max32(a.Y, b.Y)}
}

func (a Vec2D) Min(b Vec2D) Vec2D {
	return Vec2D{Min32(a.X, b.X), Min32(a.Y, b.Y)}
}

// minimum of all positive (> 0) numbers
func (a Vec2D) MinPos(b Vec2D) Vec2D {
	return Vec2D{MinPos32(a.X, b.X), MinPos32(a.Y, b.Y)}
}

// set to max of current vs. b
func (a *Vec2D) SetMax(b Vec2D) {
	a.X = Max32(a.X, b.X)
	a.Y = Max32(a.Y, b.Y)
}

// set to min of current vs. b
func (a *Vec2D) SetMin(b Vec2D) {
	a.X = Min32(a.X, b.X)
	a.Y = Min32(a.Y, b.Y)
}

// set to minpos of current vs. b
func (a *Vec2D) SetMinPos(b Vec2D) {
	a.X = MinPos32(a.X, b.X)
	a.Y = MinPos32(a.Y, b.Y)
}

// set to max of current value and val
func (a *Vec2D) SetMaxVal(val float32) {
	a.X = Max32(a.X, val)
	a.Y = Max32(a.Y, val)
}

// set to min of current value and val
func (a *Vec2D) SetMinVal(val float32) {
	a.X = Min32(a.X, val)
	a.Y = Min32(a.Y, val)
}

// set to minpos of current value and val
func (a *Vec2D) SetMinPosVal(val float32) {
	a.X = MinPos32(a.X, val)
	a.Y = MinPos32(a.Y, val)
}

// set the value along a given dimension to max of current val and new val
func (a *Vec2D) SetMaxDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X = Max32(a.X, val)
	case Y:
		a.Y = Max32(a.Y, val)
	}
}

// set the value along a given dimension to min of current val and new val
func (a *Vec2D) SetMinDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X = Min32(a.X, val)
	case Y:
		a.Y = Min32(a.Y, val)
	}
}

// set the value along a given dimension to min of current val and new val
func (a *Vec2D) SetMinPosDim(d Dims2D, val float32) {
	switch d {
	case X:
		a.X = MinPos32(val, a.X)
	case Y:
		a.Y = MinPos32(val, a.Y)
	}
}

func (a Vec2D) Abs() Vec2D {
	b := a
	if b.X < 0 {
		b.X = -b.X
	}
	if b.Y < 0 {
		b.Y = -b.Y
	}
	return b
}

func (a *Vec2D) SetPoint(pt image.Point) {
	a.X = float32(pt.X)
	a.Y = float32(pt.Y)
}

func (a *Vec2D) SetFixed(pt fixed.Point26_6) {
	a.X = FixedToFloat32(pt.X)
	a.Y = FixedToFloat32(pt.Y)
}

func (a Vec2D) ToPoint() image.Point {
	return image.Point{int(a.X), int(a.Y)}
}

func (a Vec2D) ToPointCeil() image.Point {
	return image.Point{int(math32.Ceil(a.X)), int(math32.Ceil(a.Y))}
}

func (a Vec2D) ToPointFloor() image.Point {
	return image.Point{int(math32.Floor(a.X)), int(math32.Floor(a.Y))}
}

func (a Vec2D) ToPointRound() image.Point {
	return image.Point{int(math.Round(float64(a.X))), int(math.Round(float64(a.Y)))}
}

// RectFromPosSizeMax returns an image.Rectangle from max dims of pos, size
// (floor on pos, ceil on size)
func RectFromPosSizeMax(pos, sz Vec2D) image.Rectangle {
	tp := pos.ToPointFloor()
	ts := sz.ToPointCeil()
	return image.Rect(tp.X, tp.Y, tp.X+ts.X, tp.Y+ts.Y)
}

// RectFromPosSizeMin returns an image.Rectangle from min dims of pos, size
// (ceil on pos, floor on size)
func RectFromPosSizeMin(pos, sz Vec2D) image.Rectangle {
	tp := pos.ToPointCeil()
	ts := sz.ToPointFloor()
	return image.Rect(tp.X, tp.Y, tp.X+ts.X, tp.Y+ts.Y)
}

func (a Vec2D) Distance(b Vec2D) float32 {
	return math32.Hypot(a.X-b.X, a.Y-b.Y)
}

func (a Vec2D) Interpolate(b Vec2D, t float32) Vec2D {
	x := a.X + (b.X-a.X)*t
	y := a.Y + (b.Y-a.Y)*t
	return Vec2D{x, y}
}

func (a Vec2D) String() string {
	return fmt.Sprintf("(%v, %v)", a.X, a.Y)
}

////////////////////////////////////////////////////////////////////////////////////////
// Matrix2D

// todo: in theory a high-quality SVG implementation should use a 64bit xform
// matrix, but that is rather inconvenient and unlikely to be relevant here..
// revisit later

type Matrix2D struct {
	XX, YX, XY, YY, X0, Y0 float32
}

var KiT_Matrix2D = kit.Types.AddType(&Matrix2D{}, Matrix2DProps)

var Matrix2DProps = ki.Props{
	"style-prop": true,
}

func Identity2D() Matrix2D {
	return Matrix2D{
		1, 0,
		0, 1,
		0, 0,
	}
}

func Translate2D(x, y float32) Matrix2D {
	return Matrix2D{
		1, 0,
		0, 1,
		x, y,
	}
}

func Scale2D(x, y float32) Matrix2D {
	return Matrix2D{
		x, 0,
		0, y,
		0, 0,
	}
}

func Rotate2D(angle float32) Matrix2D {
	c := float32(math32.Cos(angle))
	s := float32(math32.Sin(angle))
	return Matrix2D{
		c, s,
		-s, c,
		0, 0,
	}
}

func Shear2D(x, y float32) Matrix2D {
	return Matrix2D{
		1, y,
		x, 1,
		0, 0,
	}
}

func Skew2D(x, y float32) Matrix2D {
	return Matrix2D{
		1, math32.Tan(y),
		math32.Tan(x), 1,
		0, 0,
	}
}

func (a Matrix2D) Multiply(b Matrix2D) Matrix2D {
	return Matrix2D{
		a.XX*b.XX + a.YX*b.XY,
		a.XX*b.YX + a.YX*b.YY,
		a.XY*b.XX + a.YY*b.XY,
		a.XY*b.YX + a.YY*b.YY,
		a.X0*b.XX + a.Y0*b.XY + b.X0,
		a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

func (a Matrix2D) TransformVector(x, y float32) (tx, ty float32) {
	tx = a.XX*x + a.XY*y
	ty = a.YX*x + a.YY*y
	return
}

func (a Matrix2D) TransformVectorVec2D(v Vec2D) Vec2D {
	tx := a.XX*v.X + a.XY*v.Y
	ty := a.YX*v.X + a.YY*v.Y
	return Vec2D{tx, ty}
}

func (a Matrix2D) TransformPoint(x, y float32) (tx, ty float32) {
	tx = a.XX*x + a.XY*y + a.X0
	ty = a.YX*x + a.YY*y + a.Y0
	return
}

func (a Matrix2D) TransformPointVec2D(v Vec2D) Vec2D {
	tx := a.XX*v.X + a.XY*v.Y + a.X0
	ty := a.YX*v.X + a.YY*v.Y + a.Y0
	return Vec2D{tx, ty}
}

func (a Matrix2D) TransformPointToInt(x, y float32) (tx, ty int) {
	tx = int(a.XX*x + a.XY*y + a.X0)
	ty = int(a.YX*x + a.YY*y + a.Y0)
	return
}

func (a Matrix2D) Translate(x, y float32) Matrix2D {
	return Translate2D(x, y).Multiply(a)
}

func (a Matrix2D) Scale(x, y float32) Matrix2D {
	return Scale2D(x, y).Multiply(a)
}

func (a Matrix2D) Rotate(angle float32) Matrix2D {
	return Rotate2D(angle).Multiply(a)
}

func (a Matrix2D) Shear(x, y float32) Matrix2D {
	return Shear2D(x, y).Multiply(a)
}

func (a Matrix2D) Skew(x, y float32) Matrix2D {
	return Skew2D(x, y).Multiply(a)
}

func (a Matrix2D) ToRasterx() rasterx.Matrix2D {
	return rasterx.Matrix2D{float64(a.XX), float64(a.YX), float64(a.XY), float64(a.YY), float64(a.X0), float64(a.Y0)}
}

// ExtractRot extracts the rotation component from a given matrix
func (a Matrix2D) ExtractRot() float32 {
	return math32.Atan2(-a.XY, a.XX)
}

// ExtractXYScale extracts the X and Y scale factors after undoing any
// rotation present -- i.e., in the original X, Y coordinates
func (a Matrix2D) ExtractScale() (scx, scy float32) {
	rot := a.ExtractRot()
	tx := a.Rotate(-rot)
	scx, _ = tx.TransformVector(1, 0)
	_, scy = tx.TransformVector(0, 1)
	return
}

// ParseFloat32 logs any strconv.ParseFloat errors
func ParseFloat32(pstr string) (float32, error) {
	r, err := strconv.ParseFloat(pstr, 32)
	if err != nil {
		log.Printf("gi.ParseFloat32: error parsing float32 number from: %v, %v\n", pstr, err)
		return float32(0.0), err
	}
	return float32(r), nil
}

// ParseAngle32 returns radians angle from string that can specify units (deg,
// grad, rad) -- deg is assumed if not specified
func ParseAngle32(pstr string) (float32, error) {
	units := "deg"
	lstr := strings.ToLower(pstr)
	if strings.Contains(lstr, "deg") {
		units = "deg"
		lstr = strings.TrimSuffix(lstr, "deg")
	} else if strings.Contains(lstr, "grad") {
		units = "grad"
		lstr = strings.TrimSuffix(lstr, "grad")
	} else if strings.Contains(lstr, "rad") {
		units = "rad"
		lstr = strings.TrimSuffix(lstr, "rad")
	}
	r, err := strconv.ParseFloat(lstr, 32)
	if err != nil {
		log.Printf("gi.ParseAngle32: error parsing float32 number from: %v, %v\n", lstr, err)
		return float32(0.0), err
	}
	switch units {
	case "deg":
		return float32(r) * math32.Pi / 180, nil
	case "grad":
		return float32(r) * math32.Pi / 200, nil
	case "rad":
		return float32(r), nil
	}
	return float32(r), nil
}

// ReadPoints reads a set of floating point values from a SVG format number
// string -- returns a slice or nil if there was an error
func ReadPoints(pstr string) []float32 {
	lastIdx := -1
	var pts []float32
	lr := ' '
	for i, r := range pstr {
		if unicode.IsNumber(r) == false && r != '.' && !(r == '-' && lr == 'e') && r != 'e' {
			if lastIdx != -1 {
				s := pstr[lastIdx:i]
				p, err := ParseFloat32(s)
				if err != nil {
					return nil
				}
				pts = append(pts, p)
			}
			if r == '-' {
				lastIdx = i
			} else {
				lastIdx = -1
			}
		} else if lastIdx == -1 {
			lastIdx = i
		}
		lr = r
	}
	if lastIdx != -1 && lastIdx != len(pstr) {
		s := pstr[lastIdx:len(pstr)]
		p, err := ParseFloat32(s)
		if err != nil {
			return nil
		}
		pts = append(pts, p)
	}
	return pts
}

// PointsCheckN checks the number of points read and emits an error if not equal to n
func PointsCheckN(pts []float32, n int, errmsg string) error {
	if len(pts) != n {
		return fmt.Errorf("%v incorrect number of points: %v != %v\n", errmsg, len(pts), n)
	}
	return nil
}

// SetString processes the standard SVG-style transform strings
func (a *Matrix2D) SetString(str string) error {
	errmsg := "gi.Matrix2D SetString"
	str = strings.ToLower(strings.TrimSpace(str))
	*a = Identity2D()
	if str == "none" {
		return nil
	}
	// could have multiple transforms
	for {
		pidx := strings.IndexByte(str, '(')
		if pidx < 0 {
			err := fmt.Errorf("gi.Matrix2D SetString: no params for xform: %v\n", str)
			log.Println(err)
			return err
		}
		cmd := str[:pidx]
		vals := str[pidx+1:]
		nxt := ""
		eidx := strings.IndexByte(vals, ')')
		if eidx > 0 {
			nxt = strings.TrimSpace(vals[eidx+1:])
			if strings.HasPrefix(nxt, ";") {
				nxt = strings.TrimSpace(strings.TrimPrefix(nxt, ";"))
			}
			vals = vals[:eidx]
		}
		pts := ReadPoints(vals)
		switch cmd {
		case "matrix":
			if err := PointsCheckN(pts, 6, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = Matrix2D{pts[0], pts[1], pts[2], pts[3], pts[4], pts[5]}
		case "translate":
			if err := PointsCheckN(pts, 2, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(pts[0], pts[1])
		case "translatex":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(pts[0], 0)
		case "translatey":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(0, pts[0])
		case "scale":
			if len(pts) == 1 {
				*a = a.Scale(pts[0], pts[0])
			} else if len(pts) == 2 {
				*a = a.Scale(pts[0], pts[1])
			} else {
				err := fmt.Errorf("%v incorrect number of points: 2 != %v\n", errmsg, len(pts))
				log.Println(err)
			}
		case "scalex":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Scale(pts[0], 1)
		case "scaley":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Scale(1, pts[0])
		case "rotate":
			ang := pts[0] * math32.Pi / 180 // always in degrees in this form
			if len(pts) == 3 {
				*a = a.Translate(pts[1], pts[2]).Rotate(ang).Translate(-pts[1], -pts[2])
			} else if len(pts) == 1 {
				*a = a.Rotate(ang)
			} else {
				return PointsCheckN(pts, 1, errmsg)
			}
		case "skew":
			if err := PointsCheckN(pts, 2, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(pts[0], pts[1])
		case "skewx":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(pts[0], 0)
		case "skewy":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(0, pts[0])
		}
		if nxt == "" {
			break
		}
		if !strings.Contains(nxt, "(") {
			break
		}
		str = nxt
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Geom2DInt

// Geom2DInt defines a geometry in 2D dots units (int) -- this is just a more
// convenient format than image.Rectangle for cases where the size and
// position are independently updated (e.g., Viewport)
type Geom2DInt struct {
	Pos  image.Point
	Size image.Point
}

// Bounds converts geom to equivalent image.Rectangle
func (gm *Geom2DInt) Bounds() image.Rectangle {
	return image.Rect(gm.Pos.X, gm.Pos.Y, gm.Pos.X+gm.Size.X, gm.Pos.Y+gm.Size.Y)
}

// SizeRect converts geom to rect version of size at 0 pos
func (gm *Geom2DInt) SizeRect() image.Rectangle {
	return image.Rect(0, 0, gm.Size.X, gm.Size.Y)
}

// SetRect sets values from image.Rectangle
func (gm *Geom2DInt) SetRect(r image.Rectangle) {
	gm.Pos = r.Min
	gm.Size = r.Size()
}

///////////////////////////////////////////////////////////
// utilities

func Radians(degrees float32) float32 {
	return degrees * math32.Pi / 180
}

func Degrees(radians float32) float32 {
	return radians * 180 / math32.Pi
}

func Float32ToFixedPoint(x, y float32) fixed.Point26_6 {
	return fixed.Point26_6{X: Float32ToFixed(x), Y: Float32ToFixed(y)}
}

func Float32ToFixed(x float32) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

func FixedToFloat32(x fixed.Int26_6) float32 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float32(x>>shift) + float32(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float32(x>>shift) + float32(x&mask)/64)
	}
	return 0
}

func FixedToFloat(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}
