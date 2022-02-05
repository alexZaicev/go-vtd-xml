# go-vtd-xml

[![CircleCI Build Status](https://circleci.com/gh/CircleCI-Public/circleci-demo-go.svg?style=shield)](https://circleci.com/gh/alexZaicev/go-vtd-xml)
[![Coverage Status](https://coveralls.io/repos/github/alexZaicev/go-vtd-xml/badge.svg)](https://coveralls.io/github/alexZaicev/go-vtd-xml)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/alexZaicev/go-vtd-xml/main/LICENSE.md)

GoLang implementation of a well-known Java XML parsing library ximpleware/vtd-xml.

### TODO
1. Find alternative to Java CUP that can be implemented in GoLang to parse XPath and evaluate expression 

Nice to have:
- Make sure all errors return from functions in parent functions are wrapped with appropriate error and message
- Make sure all interfaces/structs/functions have code comments