package models

import (
	"context"
	"fmt"
	"time"

	"github.com/Nadeem-Zaidi/CRM/database"
)

type Product struct {
	ID   int
	Name string

	Brand string
	VR    []Variations
}

func (p *Product) All() []Product {
	var v Variations
	variations := make([]Variations, 0)
	product := make([]Product, 0)

	query := `select * from products`
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stmt, err := database.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand); err != nil {
			fmt.Println(err)

		}

		query2 := fmt.Sprintf("select * from variations where product=%d", p.ID)

		stmt2, err := database.DB.PrepareContext(ctx, query2)
		if err != nil {
			fmt.Println(err)
		}
		defer stmt2.Close()
		rows2, err := stmt2.QueryContext(ctx)
		if err != nil {
			fmt.Println(err)
		}
		defer rows2.Close()

		for rows2.Next() {
			if err := rows2.Scan(&v.ID, &v.Color, &v.Size, &v.Price, &v.Product); err != nil {
				fmt.Println(err)

			}
			v2 := Variations{ID: v.ID, Color: v.Color, Size: v.Size, Price: v.Price, Product: v.Product}
			variations = append(variations, v2)

		}
		p.VR = variations
		product = append(product, *p)

	}
	return product

}

type VProduct struct {
	Name  string
	Brand string
	Color string
	Size  string
	Price float64
}

func (pv *VProduct) VAll() []VProduct {
	var vp VProduct
	query := "SELECT p.name,p.brand,v.color,v.size,v.price from products as p inner join variations as v  on v.product=p.id group by v.product"

	ct, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	st, err := database.DB.PrepareContext(ct, query)
	if err != nil {
		fmt.Println(err)
	}
	defer st.Close()
	rows, err := st.QueryContext(ct)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	vpl := make([]VProduct, 0)

	for rows.Next() {
		if err := rows.Scan(&vp.Name, &vp.Brand, &vp.Color, &vp.Size, &vp.Price); err != nil {
			fmt.Println(err)

		}
		vpl = append(vpl, vp)

	}
	return vpl

}

func (vp *VProduct) SingleProduct(id int) []VProduct {
	var v VProduct

	query := fmt.Sprintf("select p.name,p.brand,v.color,v.size,v.price from products as p inner join variations v on p.id=v.product where v.product=%d", id)
	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()
	stm, err := database.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	defer stm.Close()
	row, err := stm.QueryContext(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()
	vl := make([]VProduct, 0)

	for row.Next() {
		if err := row.Scan(&v.Name, &v.Brand, &v.Color, &v.Size, &v.Price); err != nil {
			fmt.Println(err)

		}
		vl = append(vl, v)

	}
	return vl

}
