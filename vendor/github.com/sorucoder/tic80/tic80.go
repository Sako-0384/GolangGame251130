package tic80

import (
	"reflect"
	"unsafe"
)

// Memory Areas
var (
	IO_RAM   = (*[0x18000]byte)(unsafe.Pointer(uintptr(0x00000)))
	FREE_RAM = (*[0x28000]byte)(unsafe.Pointer(uintptr(0x18000)))
)

// toTextData transforms a Go string into a form useable by TIC-80.
func toTextData(goString *string) unsafe.Pointer {
	textData := new([]byte)
	*textData = make([]byte, 0, len(*goString)+1)
	for _, goRune := range *goString {
		if goRune > 0 {
			switch {
			case goRune <= 0x7F:
				*textData = append(*textData, byte(goRune))
			default:
				*textData = append(*textData, byte('?'))
			}
		}
	}
	*textData = append(*textData, 0)
	buffer, _ := toByteData(textData)
	return buffer
}

// toByteData transforms a Go slice of bytes into a form useable by TIC-80.
func toByteData(goBytes *[]byte) (buffer unsafe.Pointer, count int) {
	if goBytes != nil {
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(goBytes))
		buffer = unsafe.Pointer(sliceHeader.Data)
		// See https://tinygo.org/docs/guides/compatibility/#reflectsliceheader-and-reflectstringheader.
		count = int(sliceHeader.Len)
	}
	return
}

// paletteSet represents a subset of the color palette.
type paletteSet uint16

func (mask *paletteSet) Clear() {
	*mask = 0
}

// AddColor adds a color to the set.
func (mask *paletteSet) AddColor(color int) {
	*mask |= paletteSet(1 << (color % 16))
}

// RemoveColor removes a color from the set.
func (mask *paletteSet) RemoveColor(color int) {
	*mask &^= paletteSet(1 << (color % 16))
}

// Colors returns a slice containing the colors.
func (mask *paletteSet) Colors() []byte {
	if *mask > 0 {
		set := make([]byte, 0, 16)
		for color := 0; color < 16; color++ {
			if *mask&paletteSet(1<<color) > 0 {
				set = append(set, byte(color))
			}
		}
		return set
	}
	return nil
}

// ButtonCode represents a button id for use with [tic80.Btn] and [tic80.Btnp]
type ButtonCode int

// Gamepad Players
const (
	GAMEPAD_1 ButtonCode = 8 * iota
	GAMEPAD_2
	GAMEPAD_3
	GAMEPAD_4
)

// Gamepad Buttons
const (
	BUTTON_UP ButtonCode = iota
	BUTTON_DOWN
	BUTTON_LEFT
	BUTTON_RIGHT
	BUTTON_A
	BUTTON_B
	BUTTON_X
	BUTTON_Y
)

// KeyCode represents a keyboard id for use with [tic80.Key] and [tic80.Keyp]
type KeyCode int

// Keyboard keys.
const (
	KEY_A KeyCode = iota + 1
	KEY_B
	KEY_C
	KEY_D
	KEY_E
	KEY_F
	KEY_G
	KEY_H
	KEY_I
	KEY_J
	KEY_K
	KEY_L
	KEY_M
	KEY_N
	KEY_O
	KEY_P
	KEY_Q
	KEY_R
	KEY_S
	KEY_T
	KEY_U
	KEY_V
	KEY_W
	KEY_X
	KEY_Y
	KEY_Z
	KEY_ZERO
	KEY_ONE
	KEY_TWO
	KEY_THREE
	KEY_FOUR
	KEY_FIVE
	KEY_SIX
	KEY_SEVEN
	KEY_EIGHT
	KEY_NINE
	KEY_MINUS
	KEY_EQUALS
	KEY_LEFTBRACKET
	KEY_RIGHTBRACKET
	KEY_BACKSLASH
	KEY_SEMICOLON
	KEY_APOSTROPHE
	KEY_GRAVE
	KEY_COMMA
	KEY_PERIOD
	KEY_SLASH
	KEY_SPACE
	KEY_TAB
	KEY_RETURN
	KEY_BACKSPACE
	KEY_DELETE
	KEY_INSERT
	KEY_PAGEUP
	KEY_PAGEDOWN
	KEY_HOME
	KEY_END
	KEY_UP
	KEY_DOWN
	KEY_LEFT
	KEY_RIGHT
	KEY_CAPSLOCK
	KEY_CTRL
	KEY_SHIFT
	KEY_ALT
)

// FontOptions provides additional options to [tic80.Font].
type FontOptions struct {
	transparentColors paletteSet
	characterWidth    int
	characterHeight   int
	fixed             bool
	scale             int
	alternateFont     bool
}

var defaultFontOptions FontOptions = FontOptions{
	transparentColors: 0,
	characterWidth:    8,
	characterHeight:   8,
	fixed:             false,
	scale:             1,
	alternateFont:     false,
}

// NewFontOptions constructs a [tic80.FontOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/font
func NewFontOptions() *FontOptions {
	options := new(FontOptions)
	*options = defaultFontOptions
	return options
}

// AddTransparentColor adds an additional color to the list of colors to render as transparent.
func (options *FontOptions) AddTransparentColor(color int) *FontOptions {
	options.transparentColors.AddColor(color)
	return options
}

// RemoveTransparentColor removes a color to the list of colors to render as transparent.
func (options *FontOptions) RemoveTransparentColor(color int) *FontOptions {
	options.transparentColors.RemoveColor(color)
	return options
}

// SetOpaque removes all colors to render transparent.
func (options *FontOptions) SetOpaque() *FontOptions {
	options.transparentColors.Clear()
	return options
}

// SetCharacterSize sets the maximum size of each character in pixels.
func (options *FontOptions) SetCharacterSize(width, height int) *FontOptions {
	options.characterWidth = width
	options.characterHeight = height
	return options
}

// SetScale sets the scale as a whole-number multiplier.
func (options *FontOptions) SetScale(scale int) *FontOptions {
	options.scale = scale
	return options
}

// ToggleFixed toggles monospacing.
func (options *FontOptions) ToggleFixed() *FontOptions {
	options.fixed = !options.fixed
	return options
}

// TogglePage toggles which font page to use (usually between large and small font).
func (options *FontOptions) TogglePage() *FontOptions {
	options.alternateFont = !options.alternateFont
	return options
}

// MapOptions provides additional options to [tic80.Map].
type MapOptions struct {
	x                 int
	y                 int
	width             int
	height            int
	screenX           int
	screenY           int
	transparentColors paletteSet
	scale             int
}

var defaultMapOptions MapOptions = MapOptions{
	x:                 0,
	y:                 0,
	width:             30,
	height:            17,
	screenX:           0,
	screenY:           0,
	transparentColors: 0,
	scale:             1,
}

// NewMapOptions constructs a [tic80.MapOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/map
func NewMapOptions() *MapOptions {
	options := new(MapOptions)
	*options = defaultMapOptions
	return options
}

// AddTransparentColor adds an additional color to the list of colors to render as transparent.
func (options *MapOptions) AddTransparentColor(color int) *MapOptions {
	options.transparentColors.AddColor(color)
	return options
}

// RemoveTransparentColor removes a color to the list of colors to render as transparent.
func (options *MapOptions) RemoveTransparentColor(color int) *MapOptions {
	options.transparentColors.RemoveColor(color)
	return options
}

// SetOpaque removes all colors to render transparent.
func (options *MapOptions) SetOpaque() *MapOptions {
	options.transparentColors.Clear()
	return options
}

// SetOffset sets the map coordinates in which to start drawing the map.
func (options *MapOptions) SetOffset(x, y int) *MapOptions {
	options.x = x
	options.y = y
	return options
}

// SetSize sets the size of the map to draw.
func (options *MapOptions) SetSize(width, height int) *MapOptions {
	options.width = width
	options.height = height
	return options
}

// SetPosition sets the screen coordinates to draw the map to.
func (options *MapOptions) SetPosition(x, y int) *MapOptions {
	options.screenX = x
	options.screenY = y
	return options
}

// SetScale sets the scale as a whole-number multiplier.
func (options *MapOptions) SetScale(scale int) *MapOptions {
	options.scale = scale
	return options
}

// MusicOptions provides additional options to [tic80.Music]
type MusicOptions struct {
	track   int
	frame   int
	row     int
	loop    bool
	sustain bool
	tempo   int
	speed   int
}

var defaultMusicOptions MusicOptions = MusicOptions{
	track:   -1,
	frame:   -1,
	row:     -1,
	loop:    true,
	sustain: false,
	tempo:   -1,
	speed:   -1,
}

// NewMusicOptions constructs a [tic80.MusicOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/music
func NewMusicOptions() *MusicOptions {
	options := new(MusicOptions)
	*options = defaultMusicOptions
	return options
}

// SetTrack sets the track index to start playing.
func (options *MusicOptions) SetTrack(track int) *MusicOptions {
	options.track = track % 8
	return options
}

// SetFrame sets the frame index to start playing.
func (options *MusicOptions) SetFrame(frame int) *MusicOptions {
	options.frame = frame % 16
	return options
}

// SetRow sets the row index to start playing.
func (options *MusicOptions) SetRow(row int) *MusicOptions {
	options.row = row % 64
	return options
}

// SetTempo sets the tempo in beats per minute.
func (options *MusicOptions) SetTempo(tempo int) *MusicOptions {
	options.tempo = tempo%241 + 40
	return options
}

// SetSpeed sets the speed.
func (options *MusicOptions) SetSpeed(speed int) *MusicOptions {
	options.speed = speed%31 + 1
	return options
}

// ToggleLooping toggles whether to loop the track.
func (options *MusicOptions) ToggleLooping() *MusicOptions {
	options.loop = !options.loop
	return options
}

// ToggleSustain toggles whether to sustain notes or not.
func (options *MusicOptions) ToggleSustain() *MusicOptions {
	options.sustain = !options.sustain
	return options
}

// PrintOptions provides additional options to [tic80.Print]
type PrintOptions struct {
	color         byte
	fixed         bool
	scale         int
	alternateFont bool
}

var defaultPrintOptions PrintOptions = PrintOptions{
	color:         15,
	fixed:         false,
	scale:         1,
	alternateFont: false,
}

// NewPrintOptions constructs a [tic80.PrintOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/print
func NewPrintOptions() *PrintOptions {
	options := new(PrintOptions)
	*options = defaultPrintOptions
	return options
}

// SetColor sets the color of the text to print.
func (options *PrintOptions) SetColor(color int) *PrintOptions {
	options.color = byte(color % 16)
	return options
}

// SetScale sets the scale as a whole-number multiplier.
func (options *PrintOptions) SetScale(scale int) *PrintOptions {
	options.scale = scale
	return options
}

// ToggleFixed toggles monospacing.
func (options *PrintOptions) ToggleFixed() *PrintOptions {
	options.fixed = !options.fixed
	return options
}

// TogglePage toggles which font page to use (usually between large and small font).
func (options *PrintOptions) TogglePage() *PrintOptions {
	options.alternateFont = !options.alternateFont
	return options
}

// SoundEffectNote is an enumeration of music notes.
type SoundEffectNote int

const NOTE_NONE SoundEffectNote = -1

// Notes
const (
	NOTE_C SoundEffectNote = iota
	NOTE_C_SHARP
	NOTE_D
	NOTE_D_SHARP
	NOTE_E
	NOTE_F
	NOTE_F_SHARP
	NOTE_G
	NOTE_G_SHARP
	NOTE_A
	NOTE_A_SHARP
	NOTE_B
)

// SoundEffectOptions provides additional options for [tic80.Sfx]
type SoundEffectOptions struct {
	id          int
	note        int
	octave      int
	duration    int
	channel     int
	leftVolume  int
	rightVolume int
	speed       int
}

var defaultSoundEffectOptions SoundEffectOptions = SoundEffectOptions{
	id:          -1,
	note:        -1,
	octave:      -1,
	duration:    -1,
	channel:     0,
	leftVolume:  15,
	rightVolume: 15,
	speed:       0,
}

// NewSoundEffectOptions constructs a [tic80.SoundEffectOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/sfx
func NewSoundEffectOptions() *SoundEffectOptions {
	options := new(SoundEffectOptions)
	*options = defaultSoundEffectOptions
	return options
}

// SetId sets the id of the sound effect to play.
func (options *SoundEffectOptions) SetId(id int) *SoundEffectOptions {
	options.id = id % 64
	return options
}

func (options *SoundEffectOptions) SetNote(note int) *SoundEffectOptions {
	options.note = note % 12
	options.octave = note / 12
	return options
}

// SetDuration sets the duration in frames to play the sound effect.
func (options *SoundEffectOptions) SetDuration(duration int) *SoundEffectOptions {
	options.duration = duration
	return options
}

// SetChannel sets the channel index to play the sound effect in.
func (options *SoundEffectOptions) SetChannel(channel int) *SoundEffectOptions {
	options.channel = channel % 4
	return options
}

// SetSpeed sets the speed of the sound effect.
func (options *SoundEffectOptions) SetSpeed(speed int) *SoundEffectOptions {
	if speed < -4 {
		options.speed = -4
	} else if speed > 3 {
		options.speed = 3
	} else {
		options.speed = speed
	}
	return options
}

// SetVolume sets the volume of both left and right speakers to the same level.
func (options *SoundEffectOptions) SetVolume(level int) *SoundEffectOptions {
	level %= 16
	options.leftVolume = level
	options.rightVolume = level
	return options
}

// SetStereoVolume sets the volume of left and right speakers independently.
func (options *SoundEffectOptions) SetStereoVolume(leftLevel, rightLevel int) *SoundEffectOptions {
	options.leftVolume = leftLevel % 16
	options.rightVolume = rightLevel % 16
	return options
}

// SpriteOptions provides additional options to [tic80.Spr]
type SpriteOptions struct {
	transparentColors paletteSet
	scale             int
	flip              int
	rotate            int
	width             int
	height            int
}

var defaultSpriteOptions SpriteOptions = SpriteOptions{
	transparentColors: 0,
	scale:             1,
	flip:              0,
	rotate:            0,
	width:             1,
	height:            1,
}

// NewSpriteOptions constructs a [tic80.SpriteOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/spr
func NewSpriteOptions() *SpriteOptions {
	options := new(SpriteOptions)
	*options = defaultSpriteOptions
	return options
}

// AddTransparentColor adds an additional color to the list of colors to render as transparent.
func (options *SpriteOptions) AddTransparentColor(color int) *SpriteOptions {
	options.transparentColors.AddColor(color)
	return options
}

// RemoveTransparentColor removes a color to the list of colors to render as transparent.
func (options *SpriteOptions) RemoveTransparentColor(color int) *SpriteOptions {
	options.transparentColors.RemoveColor(color)
	return options
}

// SetOpaque removes all colors to render transparent.
func (options *SpriteOptions) SetOpaque() *SpriteOptions {
	options.transparentColors.Clear()
	return options
}

// SetScale sets the scale as a whole-number multiplier.
func (options *SpriteOptions) SetScale(scale int) *SpriteOptions {
	options.scale = scale
	return options
}

// FlipHorizontally toggles horizontally flipping.
func (options *SpriteOptions) FlipHorizontally() *SpriteOptions {
	options.flip ^= 1
	return options
}

// FlipVertically toggles vertical flipping.
func (options *SpriteOptions) FlipVertically() *SpriteOptions {
	options.flip ^= 2
	return options
}

// Rotate90CW rotates the sprite 90 degrees clockwise.
func (options *SpriteOptions) Rotate90CW() *SpriteOptions {
	options.rotate = (options.rotate + 1) % 4
	return options
}

// Rotate90CCW rotates the sprite 90 degrees counterclockwise.
func (options *SpriteOptions) Rotate90CCW() *SpriteOptions {
	options.rotate = (options.rotate - 1) % 4
	return options
}

// Rotate180 rotates the sprite 180 degrees.
func (options *SpriteOptions) Rotate180() *SpriteOptions {
	options.rotate = (options.rotate + 2) % 4
	return options
}

// SetSize sets the size of the sprite in 8x8 sub-sprites.
func (options *SpriteOptions) SetSize(width, height int) *SpriteOptions {
	options.width = width
	options.height = height
	return options
}

// SyncMask is an enumeration of data banks.
type SyncMask int

const SYNC_ALL SyncMask = 0

// Banks
const (
	SYNC_TILES SyncMask = 1 << iota
	SYNC_SPRITES
	SYNC_MAP
	SYNC_SOUND_EFFECTS
	SYNC_MUSIC
	SYNC_PALETTE
	SYNC_FLAGS
	SYNC_SCREEN
)

// TexturedTriangleOptions provides additional options for [tic80.Ttri]
type TexturedTriangleOptions struct {
	useTiles             bool
	transparentColors    paletteSet
	useDepthCalculations bool
	z0                   int
	z1                   int
	z2                   int
}

var defaultTexturedTriangleOptions TexturedTriangleOptions = TexturedTriangleOptions{
	useTiles:             false,
	transparentColors:    0,
	useDepthCalculations: false,
	z0:                   0,
	z1:                   0,
	z2:                   0,
}

// NewTexturedTriangleOptions constructs a [tic80.TexturedTriangleOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/ttri
func NewTexturedTriangleOptions() *TexturedTriangleOptions {
	options := new(TexturedTriangleOptions)
	*options = defaultTexturedTriangleOptions
	return options
}

// AddTransparentColor adds an additional color to the list of colors to render as transparent.
func (options *TexturedTriangleOptions) AddTransparentColor(color int) *TexturedTriangleOptions {
	options.transparentColors.AddColor(color)
	return options
}

// RemoveTransparentColor removes a color to the list of colors to render as transparent.
func (options *TexturedTriangleOptions) RemoveTransparentColor(color int) *TexturedTriangleOptions {
	options.transparentColors.RemoveColor(color)
	return options
}

// SetOpaque removes all colors to render transparent.
func (options *TexturedTriangleOptions) SetOpaque() *TexturedTriangleOptions {
	options.transparentColors.Clear()
	return options
}

// SetTextureDepth enables z-depth consideration and sets the z-depth for each vertex of the triangle.
func (options *TexturedTriangleOptions) SetTextureDepth(z0, z1, z2 int) *TexturedTriangleOptions {
	options.useDepthCalculations = true
	options.z0 = z0
	options.z1 = z1
	options.z2 = z2
	return options
}

// ToggleTextureSource toggles whether to use tiles or sprites for the texture source.
func (options *TexturedTriangleOptions) ToggleTextureSource() *TexturedTriangleOptions {
	options.useTiles = !options.useTiles
	return options
}

// TraceOptions provides additional options for [tic80.Trace]
type TraceOptions struct {
	color byte
}

var defaultTraceOptions = TraceOptions{
	color: 15,
}

// NewTraceOptions constructs a [tic80.TraceOptions] object with the defaults.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/trace
func NewTraceOptions() *TraceOptions {
	options := new(TraceOptions)
	*options = defaultTraceOptions
	return options
}

// SetColor sets the color of the trace output.
func (options *TraceOptions) SetColor(color int) *TraceOptions {
	options.color = byte(color % 16)
	return options
}

//go:export btn
func rawBtn(id int32) int32

// Btn returns true if the controller button specified by the given id is pressed; false otherwise.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/btn
func Btn(id ButtonCode) bool {
	return rawBtn(int32(id%32)) > 0
}

//go:export btnp
func rawBtnp(id, hold, period int32) bool

// Btnp returns true if the controller button specified by the given id was pressed the last frame, or after hold every period frames; false otherwise.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/btnp
func Btnp(id ButtonCode, hold, period int) bool {
	return rawBtnp(int32(id%32), int32(hold), int32(period))
}

//go:export clip
func rawClip(x, y, width, height int32)

// Clip sets the clipping region for the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/clip
func Clip(x, y, width, height int) {
	rawClip(int32(x), int32(y), int32(width), int32(height))
}

//go:export cls
func rawCls(color int8)

// Cls fills the screen with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/cls
func Cls(color int) {
	rawCls(int8(color))
}

//go:export circ
func rawCirc(x, y, radius int32, color int8)

// Circ draws a filled circle with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/circ
func Circ(x, y, radius, color int) {
	rawCirc(int32(x), int32(y), int32(radius), int8(color%16))
}

//go:export circb
func rawCircb(x, y, radius int32, color int8)

// Circb draws a circle border with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/circb
func Circb(x, y, radius, color int) {
	rawCircb(int32(x), int32(y), int32(radius), int8(color%16))
}

//go:export elli
func rawElli(x, y, radiusX, radiusY int32, color int8)

// Elli draws a filled ellipse with the specified color to the screen.
func Elli(x, y, radiusX, radiusY, color int) {
	rawElli(int32(x), int32(y), int32(radiusX), int32(radiusY), int8(color%16))
}

//go:export ellib
func rawEllib(x, y, radiusX, radiusY int32, color int8)

// Ellib draws an ellipse border with the specified color to the screen.
func Ellib(x, y, radiusX, radiusY, color int) {
	rawEllib(int32(x), int32(y), int32(radiusX), int32(radiusY), int8(color%16))
}

//go:export exit
func rawExit()

// Exit closes TIC-80.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/exit
func Exit() {
	rawExit()
}

//go:export fget
func rawFget(sprite int32, flag int8) bool

// Fget gets the status of the specified flag of the specified sprite.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/fget
func Fget(sprite, flag int) bool {
	return rawFget(int32(sprite%512), int8(flag%8))
}

//go:export fset
func rawFset(sprite int32, flag int8, value bool)

// Fset sets the status of the specified flag of the specified sprite.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/fset
func Fset(sprite, flag int, value bool) {
	rawFset(int32(sprite%512), int8(flag%8), value)
}

//go:export font
func rawFont(textBuffer unsafe.Pointer, x, y int32, transparentColorBuffer unsafe.Pointer, transparentColorCount int8, characterWidth, characterHeight int8, fixed bool, scale int8, useAlternateFontPage bool) int32

// Font draws text to the screen using sprite data.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/font
func Font(text string, x, y int, options *FontOptions) (textWidth int) {
	if options == nil {
		options = &defaultFontOptions
	}

	transparentColors := options.transparentColors.Colors()
	transparentColorBuffer, transparentColorCount := toByteData(&transparentColors)
	textBuffer := toTextData(&text)

	return int(rawFont(textBuffer, int32(x), int32(y), transparentColorBuffer, int8(transparentColorCount), int8(options.characterWidth), int8(options.characterHeight), options.fixed, int8(options.scale), options.alternateFont))
}

//go:export key
func rawKey(id int32) int32

// Key returns true if keyboard key specified by the id was pressed; false otherwise.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/key
func Key(id KeyCode) bool {
	return rawKey(int32(id)) > 0
}

//go:export keyp
func rawKeyp(id int8, hold, period int32) int32

// Keyp returns true if the keyboard key specified by the given id was pressed the last frame, or after hold every period frames; false otherwise.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/btnp
func Keyp(id KeyCode, hold, period int) bool {
	return rawKeyp(int8(id), int32(hold), int32(period)) > 0
}

//go:export line
func rawLine(x0, y0, x1, y1 float32, color int8)

// Line draws a line with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/line
func Line(x0, y0, x1, y1, color int) {
	rawLine(float32(x0), float32(y0), float32(x1), float32(y1), int8(color))
}

//go:export map
func rawMap(x, y, width, height, screenX, screenY int32, transparentColorBuffer unsafe.Pointer, transparentColorCount int8, scale int32, unused int32)

// Map draws a tile map to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/map
func Map(options *MapOptions) {
	if options == nil {
		options = &defaultMapOptions
	}

	transparentColors := options.transparentColors.Colors()
	transparentColorBuffer, transparentColorCount := toByteData(&transparentColors)

	rawMap(int32(options.x), int32(options.y), int32(options.width), int32(options.height), int32(options.screenX), int32(options.screenY), transparentColorBuffer, int8(transparentColorCount), 1, 0)
}

//go:export memcpy
func rawMemcpy(destination, source, length int32)

// Memcpy copies a buffer of RAM to RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/memcpy
func Memcpy(destination, source, length int) {
	rawMemcpy(int32(destination), int32(source), int32(length))
}

//go:export memset
func rawMemset(address, value, length int32)

// Memset sets a buffer of RAM to one value.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/memset
func Memset(address, value, length int) {
	rawMemset(int32(address), int32(value), int32(length))
}

//go:export mget
func rawMget(x, y int32) int32

// Mget gets the id of a tile given by the specified coordinates on the map.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/mget
func Mget(x, y int) int {
	return int(rawMget(int32(x), int32(y)))
}

//go:export mset
func rawMset(x, y, value int32)

// Mset sets the specified id of a tile given by the specified coordinates on the map.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/mset
func Mset(x, y, value int) {
	rawMset(int32(x), int32(y), int32(value))
}

type mouseData struct {
	x       int16
	y       int16
	scrollX int8
	scrollY int8
	left    bool
	middle  bool
	right   bool
}

var mouse *mouseData = new(mouseData)

//go:export mouse
func rawMouse(data *mouseData)

// Mouse returns the current state of the mouse.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/mouse
func Mouse() (x, y int, left, middle, right bool, scrollX, scrollY int) {
	rawMouse(mouse)

	x = int(mouse.x)
	y = int(mouse.y)
	left = mouse.left
	middle = mouse.middle
	right = mouse.right
	scrollX = int(mouse.scrollX)
	scrollY = int(mouse.scrollY)
	return
}

//go:export music
func rawMusic(track, frame, row int32, loop, sustain bool, tempo, speed int32)

// Music plays a music track.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/music
func Music(options *MusicOptions) {
	if options == nil {
		options = &defaultMusicOptions
	}

	rawMusic(int32(options.track), int32(options.frame), int32(options.row), options.loop, options.sustain, int32(options.tempo), int32(options.speed))
}

//go:export peek
func rawPeek(address int32, bits int8) int8

// Peek reads a byte from RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/peek
func Peek(address int) byte {
	return byte(rawPeek(int32(address), 8))
}

// Peek4 reads a nybble from RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/peek
func Peek4(address int) byte {
	return byte(rawPeek(int32(address), 4))
}

// Peek2 reads a half-nybble from RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/peek
func Peek2(address int) byte {
	return byte(rawPeek(int32(address), 2))
}

// Peek reads a bit from RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/peek
func Peek1(address int) byte {
	return byte(rawPeek(int32(address), 1))
}

//go:export pix
func rawPix(x, y int32, color int8) uint8

// Pix draws a pixel to the screen, and returns the original color.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/pix
func Pix(x, y, color int) int {
	return int(rawPix(int32(x), int32(y), int8(color%16)))
}

//go:export pmem
func rawPmem(address int32, value int64) uint32

// Pmem reads and writes values to persistent memory.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/pmem
func Pmem(address int, value int64) uint32 {
	return rawPmem(int32(address), value)
}

//go:export poke
func rawPoke(address int32, value, bits int8)

// Poke writes a byte to RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/poke
func Poke(address int, value byte) {
	rawPoke(int32(address), int8(value), 8)
}

// Poke4 writes a nybble to RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/poke
func Poke4(address int, value byte) {
	rawPoke(int32(address), int8(value), 4)
}

// Poke2 writes a half-nybble to RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/poke
func Poke2(address int, value byte) {
	rawPoke(int32(address), int8(value), 2)
}

// Poke1 writes a bit to RAM.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/poke
func Poke1(address int, value byte) {
	rawPoke(int32(address), int8(value), 1)
}

//go:export print
func rawPrint(textBuffer unsafe.Pointer, x, y int32, color, fixed, scale, useAlternateFontPage int8) int32

// Print prints text to the screen using the system fonts.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/print
func Print(text string, x, y int, options *PrintOptions) int {
	if options == nil {
		options = &defaultPrintOptions
	}

	textBuffer := toTextData(&text)

	var optionFixed int8
	if options.fixed {
		optionFixed = 1
	}

	var optionAlternateFont int8
	if options.alternateFont {
		optionAlternateFont = 1
	}

	return int(rawPrint(textBuffer, int32(x), int32(y), int8(options.color), optionFixed, int8(options.scale), optionAlternateFont))
}

//go:export rect
func rawRect(x, y, width, height int32, color int8)

// Rect draws a filled rectangle with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/rect
func Rect(x, y, width, height, color int) {
	rawRect(int32(x), int32(y), int32(width), int32(height), int8(color%16))
}

//go:export rectb
func rawRectb(x, y, width, height int32, color int8)

// Rectb draws a rectangle border with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/rectb
func Rectb(x, y, width, height, color int) {
	rawRectb(int32(x), int32(y), int32(width), int32(height), int8(color%16))
}

//go:export sfx
func rawSfx(id, note, octave, duration, channel, volumeLeft, volumeRight, speed int32)

// Sfx plays a sound effect.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/sfx
func Sfx(options *SoundEffectOptions) {
	if options == nil {
		options = &defaultSoundEffectOptions
	}

	rawSfx(int32(options.id), int32(options.note), int32(options.octave), int32(options.duration), int32(options.channel), int32(options.leftVolume), int32(options.rightVolume), int32(options.speed))
}

//go:export spr
func rawSpr(id, x, y int32, transparentColorBuffer unsafe.Pointer, transparentColorCount int8, scale, flip, rotate, width, height int32)

// Spr draws a sprite to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/spr
func Spr(id, x, y int, options *SpriteOptions) {
	if options == nil {
		options = &defaultSpriteOptions
	}

	transparentColors := options.transparentColors.Colors()
	transparentColorBuffer, transparentColorCount := toByteData(&transparentColors)

	rawSpr(int32(id), int32(x), int32(y), transparentColorBuffer, int8(transparentColorCount), int32(options.scale), int32(options.flip), int32(options.rotate), int32(options.width), int32(options.height))
}

//go:export sync
func rawSync(mask int32, bank, toCart int8)

// Sync exchanges and optionally persists the changes of data banks.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/sync
func Sync(mask SyncMask, bank int, toCart bool) {
	var toCartValue int8
	if toCart {
		toCartValue = 1
	}

	rawSync(int32(mask), int8(bank), toCartValue)
}

//go:export ttri
func rawTtri(x0, y0, x1, y1, x2, y2, u0, v0, u1, v1, u2, v2 float32, useTiles int32, transparentColorBuffer unsafe.Pointer, transparentColorCount int8, z0, z1, z2 float32, depth bool)

// Ttri draws a textured triangle using sprites or tiles as its texture to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/ttri
func Ttri(x0, y0, x1, y1, x2, y2, u0, v0, u1, v1, u2, v2 int, options *TexturedTriangleOptions) {
	if options == nil {
		options = &defaultTexturedTriangleOptions
	}

	transparentColors := options.transparentColors.Colors()
	transparentColorBuffer, transparentColorCount := toByteData(&transparentColors)

	var useTilesValue int32
	if options.useTiles {
		useTilesValue = 1
	}

	rawTtri(float32(x0), float32(y0), float32(x1), float32(y1), float32(x2), float32(y2), float32(u0), float32(v0), float32(u1), float32(v1), float32(u2), float32(v2), useTilesValue, transparentColorBuffer, int8(transparentColorCount), float32(options.z0), float32(options.z1), float32(options.z2), options.useDepthCalculations)
}

//go:export time
func rawTime() float32

// Time returns the number of milliseconds since the game started.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/time
func Time() float32 {
	return rawTime()
}

//go:export trace
func rawTrace(messageBuffer unsafe.Pointer, color int8)

// Trace writes text to the console.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/trace
func Trace(message string, options *TraceOptions) {
	if options == nil {
		options = &defaultTraceOptions
	}

	messageBuffer := toTextData(&message)

	rawTrace(messageBuffer, int8(options.color))
}

//go:export tri
func rawTri(x0, y0, x1, y1, x2, y2 float32, color int8)

// Tri draws a filled triangle with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/tri
func Tri(x0, y0, x1, y1, x2, y2, color int) {
	rawTri(float32(x0), float32(y0), float32(x1), float32(y1), float32(x2), float32(y2), int8(color))
}

//go:export trib
func rawTrib(x0, y0, x1, y1, x2, y2 float32, color int8)

// Trib draws a triangle border with the specified color to the screen.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/trib
func Trib(x0, y0, x1, y1, x2, y2, color int) {
	rawTrib(float32(x0), float32(y0), float32(x1), float32(y1), float32(x2), float32(y2), int8(color))
}

//go:export tstamp
func rawTstamp() uint32

// Tstamp returns the current Unix timestamp.
// See the [API] for more details.
//
// [API]: https://github.com/nesbox/TIC-80/wiki/tstamp
func Tstamp() uint32 {
	return rawTstamp()
}

// Initialize is a workaround to allow TIC-80 to run Go code.
// This should be the first function run in BOOT.
//
//go:linkname Initialize _initialize
func Initialize()

//go:export main.main
func main() {}
