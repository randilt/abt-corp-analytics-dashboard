const API_BASE = "http://localhost:8080";

interface CountryRevenue {
  country: string;
  product_name: string;
  total_revenue: number;
  transaction_count: number;
}

interface ProductFrequency {
  product_id: string;
  product_name: string;
  purchase_count: number;
  current_stock: number;
}

interface MonthlySales {
  month: string;
  sales_volume: number;
  item_count: number;
}

interface RegionRevenue {
  region: string;
  total_revenue: number;
  items_sold: number;
}

interface AnalyticsResponse {
  summary: {
    total_records: number;
    processing_time_ms: number;
    cache_hit: boolean;
    country_revenue_count: number;
    top_products_count: number;
    monthly_sales_count: number;
    top_regions_count: number;
    total_revenue: number;
  };
  country_revenue: CountryRevenue[];
  top_products: ProductFrequency[];
  monthly_sales: MonthlySales[];
  top_regions: RegionRevenue[];
}

interface StatsResponse {
  total_records: number;
  processing_time_ms: number;
  cache_hit: boolean;
  country_revenue_count: number;
  top_products_count: number;
  monthly_sales_count: number;
  top_regions_count: number;
}

interface CountryRevenuePaginatedResponse {
  data: CountryRevenue[];
  count: number;
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

interface DataResponse<T> {
  data: T[];
  count: number;
}

async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  try {
    const url = `${API_BASE}${endpoint}`;
    console.log(`Fetching: ${url}`);
    const response = await fetch(url, options);
    console.log(`Response status: ${response.status}`);

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    console.log(`Response data:`, data);
    return data;
  } catch (error) {
    console.error(`API Error for ${endpoint}:`, error);
    throw error;
  }
}

export async function getAnalytics(): Promise<AnalyticsResponse> {
  return fetchApi<AnalyticsResponse>("/api/v1/analytics");
}

export async function getStats(): Promise<StatsResponse> {
  return fetchApi<StatsResponse>("/api/v1/analytics/stats");
}

export async function getCountryRevenue(
  limit = 100,
  offset = 0
): Promise<CountryRevenuePaginatedResponse> {
  return fetchApi<CountryRevenuePaginatedResponse>(
    `/api/v1/analytics/country-revenue?limit=${limit}&offset=${offset}`
  );
}

export async function getTopProducts(): Promise<
  DataResponse<ProductFrequency>
> {
  return fetchApi<DataResponse<ProductFrequency>>(
    "/api/v1/analytics/top-products"
  );
}

export async function getMonthlySales(): Promise<DataResponse<MonthlySales>> {
  return fetchApi<DataResponse<MonthlySales>>(
    "/api/v1/analytics/monthly-sales"
  );
}

export async function getTopRegions(): Promise<DataResponse<RegionRevenue>> {
  return fetchApi<DataResponse<RegionRevenue>>("/api/v1/analytics/top-regions");
}

export async function refreshCache(): Promise<{ message: string }> {
  return fetchApi<{ message: string }>("/api/v1/analytics/refresh", {
    method: "POST",
  });
}

export async function healthCheck(): Promise<{ status: string }> {
  return fetchApi<{ status: string }>("/health");
}

export type {
  CountryRevenue,
  ProductFrequency,
  MonthlySales,
  RegionRevenue,
  AnalyticsResponse,
  StatsResponse,
  CountryRevenuePaginatedResponse,
  DataResponse,
};
