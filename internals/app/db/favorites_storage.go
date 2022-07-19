package db

import (
	"database/sql"
	"fmt"

	"dclibrary.com/internals/app/models"
)

type FavorietesStorage struct {
	DB *sql.DB
}

func (m *FavorietesStorage) AddToFavorieteBooks(favorietes *models.FavorieteBooks) (int, error) {

	stmt := `INSERT INTO dclibrary.favoriete_books (user_id, book_id) VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, favorietes.UserId, favorietes.FavoriteBookId)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> book %v has been added to Favorietes", id)

	return int(id), nil
}

func (m *FavorietesStorage) AddToFavorieteAuthors(favorietes *models.FavorieteAuthors) (int, error) {

	stmt := `INSERT INTO dclibrary.favoriete_authors (user_id, author_id) VALUES(?, ?)`

	result, err := m.DB.Exec(stmt, favorietes.UserId, favorietes.FavoriteAuthorId)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("---> Author %v has been added to Favorietes", id)

	return int(id), nil
}

func (m *FavorietesStorage) GetFavoriteBooksList(user_id int64) ([]*models.Books, error) {

	stmt := `SELECT user_id, book_id FROM dclibrary.favoriete_books WHERE user_id = ?`
	sbmb := `SELECT book_id, bookname, count, photo, stars FROM dclibrary.books WHERE book_id = ?`
	skmk := `SELECT author_id FROM dclibrary.books_authors WHERE book_id = ?`
	sdmd := `SELECT author_name, author_surname, author_patrynomic, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	var mybooks []*models.Books

	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.FavorieteBooks{}
		err = rows.Scan(&s.UserId, &s.FavoriteBookId)
		if err != nil {
			return nil, err
		}
		//have Favorietes{UserId: 0, FavorieteBookId: 0}

		bookstr, err := m.DB.Query(sbmb, s.FavoriteBookId)
		if err != nil {
			return nil, err
		}
		defer bookstr.Close()

		for bookstr.Next() {
			b := &models.Books{}
			err = rows.Scan(&b.BookId, &b.BookName, &b.Count, &b.BookPhoto, &b.Stars)
			if err != nil {
				return nil, err
			}
			//	have
			//	Books{BookId: 0,
			//	BookName: "NAME",
			//	Count: 0,
			//	BookPhoto: "/path/photo.png",
			//	Stars: 0,
			//	Authors{AuthorId: nil,
			//		AuthorPhoto: nil,
			//		AuthorStars: nil,
			//		AuthorName{
			//			Name: nil,
			//			Surname: nil,
			//			Patrynomic: nil}}}

			connection, err := m.DB.Query(skmk, b.BookId)
			if err != nil {
				return nil, err
			}

			defer connection.Close()

			for connection.Next() {
				a := &models.Authors{}

				err = connection.Scan(&a.AuthorId)
				if err != nil {
					return nil, err
				}
				//	have
				//	Books{BookId: 0,
				//	BookName: "NAME",
				//	Count: 0,
				//	BookPhoto: "/path/photo.png",
				//	Stars: 0,
				//	Authors{AuthorId: 0,
				//		AuthorPhoto: nil,
				//		AuthorStars: nil,
				//		AuthorName{
				//			Name: nil,
				//			Surname: nil,
				//			Patrynomic: nil}}}

				authors := m.DB.QueryRow(sdmd, a.AuthorId)
				err = authors.Scan(
					&a.AuthorName.Name,
					&a.AuthorName.Surname,
					&a.AuthorName.Patronymic,
					&a.AuthorPhoto,
					&a.AuthorStars,
				)
				//	have
				//	Books{BookId: 0,
				//	BookName: "NAME",
				//	Count: 0,
				//	BookPhoto: "/path/photo.png",
				//	Stars: 0,
				//	Author{AuthorId: 0,
				//		AuthorPhoto: "/path/authors/photo.png",
				//		AuthorStars: 0,
				//		AuthorName{
				//			Name: "NAME",
				//			Surname: "SURNAME",
				//			Patrynomic: "PATRYNOMIC"}}}

				if err != nil {
					return nil, err
				}

				b.Authors = append(b.Authors, *a)
			}

			mybooks = append(mybooks, b)
		}
	}
	return mybooks, nil
}

func (m *FavorietesStorage) GetFavoriteAuthorsList(user_id int64) ([]*models.Authors, error) {

	stmt := `SELECT user_id, author_id FROM dclibrary.favoriete_authors WHERE user_id = ?`
	sdmd := `SELECT author_name, author_surname, author_patrynomic, author_photo, author_stars FROM dclibrary.authors WHERE author_id = ?`

	var myauthors []*models.Authors

	rows, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		a := &models.Authors{}
		s := &models.FavorieteAuthors{}
		err = rows.Scan(&s.UserId, &s.FavoriteAuthorId)
		if err != nil {
			return nil, err
		}
		//have Favorietes{UserId: 0, FavoriteAuthorId: 0}

		bookstr := m.DB.QueryRow(sdmd, s.FavoriteAuthorId)

		bookstr.Scan(&a.AuthorName.Name, &a.AuthorName.Surname, &a.AuthorName.Patronymic, &a.AuthorPhoto, &a.AuthorStars)

		myauthors = append(myauthors, a)
	}
	return myauthors, nil
}

func (m *FavorietesStorage) DeleteFavorieteBookById(id int64) (int64, error) {
	stmt := `DELETE FROM dclibrary.favoriete_books WHERE book_id = ?`

	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (m *FavorietesStorage) DeleteFavorieteAuthorById(id int64) (int64, error) {
	stmt := `DELETE FROM dclibrary.favoriete_authors WHERE author_id = ?`

	deleted, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	res, err := deleted.LastInsertId()
	if err != nil {
		return 0, err
	}
	return res, nil
}
