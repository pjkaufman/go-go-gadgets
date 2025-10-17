//go:build unit

package epubcheck

/**
What should be tested here:
- Make sure each of the expected rules runs and makes the expected update (these will be single tests with edge cases allotted as well)
  - OPF 14: add scripted to a manifest file
	- OPF 15: remove scripted from a manifest file
	- NCX-001: fix book id discrepancy between ncx and opf files
	- RSC 5:
	  - Invalid id:
		  - Should work on opf files
			- Should work on ncx files
			- Should work on html/xhtml files
			- Make sure that it increments or decrements anything that comes after it on the same line based on the change in characters from the rule
		- Invalid attribute:
		  - Should swap to the refines syntax in the opf file and work on multiple entries
		- Empty metadata property:
		  - Should remove the empty metadata property in opf files, remove any errors for that line, and decrement subsequent line number references on errors
		-
- Make sure that multiple rules play well together
*/
