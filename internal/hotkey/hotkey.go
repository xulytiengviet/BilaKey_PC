package hotkey

type Action uint8

const (
	None Action = iota
	ToggleVietnamese
	SelectCVNSS
	SelectTelex
	SelectVNI
	CycleCandidate
)

// Resolve keeps all global shortcuts behind Ctrl+Shift and never steals
// Ctrl+Tab, Alt combinations, or Windows-key combinations from applications.
func Resolve(vk uint32, ctrl, shift, alt, win bool) Action {
	if !ctrl || !shift || alt || win {
		return None
	}
	switch vk {
	case 0x20:
		return ToggleVietnamese
	case '1':
		return SelectCVNSS
	case '2':
		return SelectTelex
	case '3':
		return SelectVNI
	case '0':
		return CycleCandidate
	default:
		return None
	}
}
