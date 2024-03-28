/*
Copyright 2022-2024 EscherCloud.
Copyright 2024 the Unikorn Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package formatter

//nolint:interfacebloat
type Formatter interface {
	// These are standard HTML types of markup.
	H1(a ...any)
	H2(a ...any)
	H3(a ...any)
	H4(a ...any)
	H5(a ...any)
	P(a string)
	Details(summary string, callback func())
	Code(lang string, title string, code string)
	Table()
	TableEnd()
	TH(a ...string)
	TD(a ...string)

	// These are more specialised mark up types e.g.
	// admonitions.
	TableOfContentsLevel(min int, max int)
	Warning(description string)
}
