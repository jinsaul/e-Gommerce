import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Product } from '../models/product';

/**
 * ProductService handles all HTTP calls to the Go backend's product endpoints.
 * providedIn: 'root' makes this a singleton — one instance shared across the entire app.
 */
@Injectable({ providedIn: 'root' })
export class ProductService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api/v1/products';

  /** GET /api/v1/products — returns all products */
  getProducts(): Observable<Product[]> {
    return this.http.get<Product[]>(this.apiUrl);
  }

  /** POST /api/v1/products — creates a new product */
  createProduct(product: Omit<Product, 'id' | 'createdAt'>): Observable<Product> {
    return this.http.post<Product>(this.apiUrl, product);
  }
}
