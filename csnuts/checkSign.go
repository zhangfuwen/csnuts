package csnuts

const authorStringLen=50
const tagStringLen=80
const titleStringLen=100
const contentStringLen=10000
const commentStringLen=200

func badAuthor(str string) bool {
    if len(str)<=authorStringLen {
        return false
    }
        return true
}
func badTitle(str string) bool {
    if len(str)<=titleStringLen {
        return false
    }
        return true
}
func badTag(str string) bool {
    if len(str)<=tagStringLen {
        return false
    }
        return true
}
func badContent(str string) bool {
    if len(str)<=contentStringLen {
        return false
    }
        return true
}
func badComment(str string) bool {
    if len(str)<=commentStringLen {
        return false
    }
        return true
}
