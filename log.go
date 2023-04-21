/*
 *
 * Copyright 2023 puzzlegrpcserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package puzzlegrpcserver

import (
	"fmt"

	"github.com/dvaumoron/puzzlelogger"
	"go.uber.org/zap"
)

type loggerWrapper struct {
	inner *zap.Logger
}

func (w loggerWrapper) Info(args ...any) {
	fmt.Fprint(puzzlelogger.InfoWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Infoln(args ...any) {
	fmt.Fprintln(puzzlelogger.InfoWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Infof(format string, args ...any) {
	fmt.Fprintf(puzzlelogger.InfoWrapper{Inner: w.inner}, format, args...)
}

func (w loggerWrapper) Warning(args ...any) {
	fmt.Fprint(puzzlelogger.WarnWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Warningln(args ...any) {
	fmt.Fprintln(puzzlelogger.WarnWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Warningf(format string, args ...any) {
	fmt.Fprintf(puzzlelogger.WarnWrapper{Inner: w.inner}, format, args...)
}

func (w loggerWrapper) Error(args ...any) {
	fmt.Fprint(puzzlelogger.ErrorWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Errorln(args ...any) {
	fmt.Fprintln(puzzlelogger.ErrorWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Errorf(format string, args ...any) {
	fmt.Fprintf(puzzlelogger.ErrorWrapper{Inner: w.inner}, format, args...)
}

func (w loggerWrapper) Fatal(args ...any) {
	fmt.Fprint(puzzlelogger.FatalWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Fatalln(args ...any) {
	fmt.Fprintln(puzzlelogger.FatalWrapper{Inner: w.inner}, args...)
}

func (w loggerWrapper) Fatalf(format string, args ...any) {
	fmt.Fprintf(puzzlelogger.FatalWrapper{Inner: w.inner}, format, args...)
}

func (w loggerWrapper) V(level int) bool {
	return true
}
