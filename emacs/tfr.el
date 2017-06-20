;;; tfr.el --- Read textfiles

;; Author: death <github.com/death>
;; Version: 1.0
;; Package-Requires: ()
;; Keywords: entertainment
;; URL: http://github.com/death/tfr

;; This file is not part of GNU Emacs.

;; Copyright (c) 2017 death

;; Permission is hereby granted, free of charge, to any person
;; obtaining a copy of this software and associated documentation
;; files (the "Software"), to deal in the Software without
;; restriction, including without limitation the rights to use, copy,
;; modify, merge, publish, distribute, sublicense, and/or sell copies
;; of the Software, and to permit persons to whom the Software is
;; furnished to do so, subject to the following conditions:

;; The above copyright notice and this permission notice shall be
;; included in all copies or substantial portions of the Software.

;; THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
;; EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
;; MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
;; NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
;; BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
;; ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
;; CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
;; SOFTWARE.

;;; Code:

(defvar tfr-program "tfr")

(defvar tfr-mode-map
  (let ((map (make-sparse-keymap)))
    (define-key map [?R] 'tfr-random)
    (define-key map [?N] 'tfr-next)
    (define-key map [?F] 'tfr-finish)
    (define-key map [?q] 'bury-buffer)
    (define-key map [?h] 'left-char)
    (define-key map [?j] 'next-line)
    (define-key map [?k] 'previous-line)
    (define-key map [?l] 'right-char)
    (define-key map [? ] 'scroll-up-command)
    (define-key map [?J] 'scroll-up-line)
    (define-key map [?K] 'scroll-down-line)
    map))

(define-derived-mode tfr-mode fundamental-mode "Textfile Reader"
  "Major mode for reading textfiles."
  (setq buffer-read-only t)
  (setq-local scroll-margin 0))

;;;###autoload
(defun tfr ()
  "Read textfiles."
  (interactive)
  (let ((buffer (get-buffer-create "*TextfileReader*")))
    (switch-to-buffer buffer)
    (tfr-mode)
    (tfr-next)))

(defun tfr-command (command &optional output)
  "Run tfr with the supplied command.

If `output' is `discard', then tfr's output will be discarded.
Otherwise, the output will be inserted into the current buffer."
  (let ((inhibit-read-only t))
    (unless (eq output 'discard)
      (erase-buffer))
    (save-excursion
      (message "TFR %s..." command)
      (call-process tfr-program
                    nil
                    (if (eq output 'discard) nil t)
                    nil
                    command))
    (message "TFR %s...done" command)
    (set-buffer-modified-p nil)))

(defun tfr-next ()
  (interactive)
  (tfr-command "next"))

(defun tfr-random ()
  (interactive)
  (tfr-command "random"))

(defun tfr-finish ()
  (interactive)
  (tfr-command "finish" 'discard))

(provide 'tfr)

;;; tfr.el ends here
