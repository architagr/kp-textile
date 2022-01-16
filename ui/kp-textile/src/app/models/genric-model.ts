export interface ErrorDetail {
	errorCode    :string 
	errorMessage :string 
}

export interface CommonResponse {
	statusCode   :number           
	errorMessage :string        
	errors       :ErrorDetail[]
}

export interface CommonListResponse extends  CommonResponse{
	lastEvalutionKey :any 
	pageSize         :number                              
	total            :number                              
}