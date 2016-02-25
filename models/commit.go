package models

import (
	"bufio"
	"bytes"
	"io"
	"time"
	"strconv"
)

/*
{
  "sha": "7638417db6d59f3c431d3e1f261cc637155684cd",
  "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/7638417db6d59f3c431d3e1f261cc637155684cd",
  "author": {
    "date": "2014-11-07T22:01:45Z",
    "name": "Scott Chacon",
    "email": "schacon@gmail.com"
  },
  "committer": {
    "date": "2014-11-07T22:01:45Z",
    "name": "Scott Chacon",
    "email": "schacon@gmail.com"
  },
  "message": "added readme, because im a good github citizen\n",
  "tree": {
    "url": "https://api.github.com/repos/octocat/Hello-World/git/trees/691272480426f78a0138979dd3ce63b77f706feb",
    "sha": "691272480426f78a0138979dd3ce63b77f706feb"
  },
  "parents": [
    {
      "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/1acc419d4d6a9ce985db7be48c6349a0475975b5",
      "sha": "1acc419d4d6a9ce985db7be48c6349a0475975b5"
    }
  ]
}
*/
type Commit struct {
	Hash      Hash
	Author    Signature
	Committer Signature
	Message   string

	tree    Hash
	parents []Hash
	r       *Repository
}

// Decode transform an Object into a Blob struct
func (c *Commit) Decode(o Object) error {
	c.Hash = o.Hash()
	r := bufio.NewReader(o.Reader())
//	fmt.Printf("######   Commit Decode r: %v ########\r\n", r)
	var message bool
	for {
		line, err := r.ReadSlice('\n')
//	    fmt.Printf("######   Commit Decode line: %v ########\r\n", string(line))
		if err != nil && err != io.EOF {
			return err
		}

		line = bytes.TrimSpace(line)
		if !message {
			if len(line) == 0 {
				message = true
				continue
			}

			split := bytes.SplitN(line, []byte{' '}, 2)
			switch string(split[0]) {
			case "tree":
				c.tree = NewHash(string(split[1]))
			case "parent":
				c.parents = append(c.parents, NewHash(string(split[1])))
			case "author":
				c.Author.Decode(split[1])
			case "committer":
				c.Committer.Decode(split[1])
			}
		} else {
			c.Message += string(line) + "\n"
		}

		if err == io.EOF {
			return nil
		}
	}
}

// Signature represents an action signed by a person
type Signature struct {
	Name  string
	Email string
	When  time.Time
}

// Decode decodes a byte slice into a signature
func (s *Signature) Decode(b []byte) {
	if len(b) == 0 {
		return
	}

	from := 0
	state := 'n' // n: name, e: email, t: timestamp, z: timezone
	for i := 0; ; i++ {
		var c byte
		var end bool
		if i < len(b) {
			c = b[i]
		} else {
			end = true
		}

		switch state {
		case 'n':
			if c == '<' || end {
				if i == 0 {
					break
				}
				s.Name = string(b[from : i-1])
				state = 'e'
				from = i + 1
			}
		case 'e':
			if c == '>' || end {
				s.Email = string(b[from:i])
				i++
				state = 't'
				from = i + 1
			}
		case 't':
			if c == ' ' || end {
				t, err := strconv.ParseInt(string(b[from:i]), 10, 64)
				if err == nil {
					loc := time.UTC
					ts := time.Unix(t, 0)
					if len(b[i:]) >= 6 {
						tl, err := time.Parse(" -0700", string(b[i:i+6]))
						if err == nil {
							loc = tl.Location()
						}
					}
					s.When = ts.In(loc)
				}
				end = true
			}
		}

		if end {
			break
		}
	}
}

type CommitIter struct {
	iter
}

type iter struct {
	ch chan Object
	r  *Repository

	IsClosed bool
}

