/*
Copyright © 2020 Alessandro Segala (@ItalyPaleAle)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package server

import "time"

type treeOperationReponse struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type metadataResponse struct {
	FileId   string     `json:"fileId"`
	Path     string     `json:"path"`
	Name     string     `json:"name"`
	Date     *time.Time `json:"date,omitempty"`
	MimeType string     `json:"mimeType,omitempty"`
	Size     int64      `json:"size,omitempty"`
}
