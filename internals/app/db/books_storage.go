package db

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"dclibrary.com/internals/app/models"
)

type BooksStorage struct {
	DB *sql.DB
}

func (m *BooksStorage) CreateNewBook(bookname, photo, authors string, count int) (int, error) {

	stmt := `INSERT INTO dclibrary.books (bookname, author_id, count, photo, stars) VALUES(?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, bookname, author_id, count, photo, 0)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> Book %v has been added to DB", id)

	return int(id), nil
}

func (m *BooksStorage) GetBooksList() ([]*models.FinalBooks, error) {

	stmt := `SELECT book_id, bookname, author_id, count, photo, stars FROM dclibrary.books`
	sdmd := `SELECT author_id, author_name, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allbooks []*models.FinalBooks

	for rows.Next() {
		s := &models.Books{}
		b := &models.FinalBooks{}
		err = rows.Scan(&s.BookId, &s.BookName, &s.Authors, &s.Count, &s.BookPhoto, &s.Stars)
		if err != nil {
			return nil, err
		}

		authors := strings.Split(s.Authors, ", ")
		for _, value := range authors {
			a := &models.Authors{}
			val, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}

			auth := m.DB.QueryRow(sdmd, val)
			err = auth.Scan(&a.AuthorId, &a.AuthorName, &a.AuthorPhoto, &a.AuthorStars)
			b.Authors = append(b.Authors, *a)
		}

		b.Book = *s
		allbooks = append(allbooks, b)
	}

	//if err = rows.Err(); err != nil {
	//	return nil, err
	//}

	return allbooks, nil
}

func (m *BooksStorage) GetBookById(id int64) (*models.FinalBooks, error) {

	stmt := `SELECT book_id, bookname, author_id, count, photo FROM dclibrary.books WHERE id = ?`
	sdmd := `SELECT author_name, author_photo FROM dclibrary.authors WHERE author_id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.FinalBooks{}

	err := row.Scan(&s.Book.BookId, &s.Book.BookName, &s.Book.Authors, &s.Book.Count, &s.Book.BookPhoto, &s.Book.Stars)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	authors := strings.Split(s.Book.Authors, ", ")
	for _, value := range authors {
		a := &models.Authors{}
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}

		auth := m.DB.QueryRow(sdmd, val)
		err = auth.Scan(&a.AuthorId, &a.AuthorName, &a.AuthorPhoto, &a.AuthorStars)
		s.Authors = append(s.Authors, *a)
	}

	/*
		authors := strings.Split(s.Author.AuthorId, ", ")

		auth := m.DB.QueryRow(sdmd, s.Author.AuthorId)
		err = auth.Scan(&s.Author.AuthorName, &s.Author.AuthorPhoto, &s.Author.AuthorStars)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, err
			}
		}
	*/

	return s, nil
}

func (m *BooksStorage) ChangeBookById(id int64, old *models.Books) (*models.Books, error) { //доделать с учётом нескольких авторов

	params := m.CheckParams()
	if search("book_name", params) {
		fmt.Println("Введите новое имя книги: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.BookName = scanner.Text()
	} else if search("author_id", params) {
		fmt.Println("Введите нового автора/авторов: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Authors = scanner.Text() // ДОДЕЛАТЬ
	} else if search("count", params) {
		fmt.Println("Введите новое количество: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.Count, _ = strconv.Atoi(scanner.Text())
	} else if search("photo", params) {
		fmt.Println("Выберите новое фото: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		old.BookPhoto = scanner.Text()
	}

	stmt := `UPDATE dclibrary.books SET book_name = ?, author_id = ?, count = ?, photo = ? WHERE book_id = ?`
	change, err := m.DB.Exec(stmt, old.BookName, old.Authors, old.Count, old.BookPhoto, id)
	if err != nil {
		return old, err
	}
	id, err = change.LastInsertId()
	if err != nil {
		return old, err
	}
	fmt.Printf("---> Book %v has been changed in DB", id)

	return old, nil
}

func (m *BooksStorage) DeleteBookById(id int64) (int, error) {
	stmt := `DELETE FROM dclibrary.books WHERE book_id = ?`
	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (m *BooksStorage) PutStarByBookId(id int64) error { //!!!

	stmt := `UPDATE dclibrary.books SET stars = ? WHERE book_id = ?`

	book, err := m.GetBookById(id)
	if err != nil {
		return err
	}

	book.Book.Stars += 1

	putstar, err := m.DB.Exec(stmt, book.Book.Stars, id)
	if err != nil {
		return err
	}

	id, err = putstar.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
