## Linnworks Golang SDK  Project
- This is a golang module
- Main technologies: golang, swagger generated code

## Task
- Create a Request Method for each API endpoint defined in the linnworks swagger spec
- Use Request and Response structs that has been generated in every request group in the `/models` folder
- Use Builder pattern for creating requests

## Testing
- Don't do tests

## Project Structure
- `/third_party` - pulled linnworks spec, do not edit
- `/{api_request_group}` - contains structs and methods for each api request
- `/{api_request_group}/models` - contains shared models that are generated from the spec, do not edit

## Coding Conventions

### Naming
- Use CamelCase for struct names and method names
- Request Methods should be named the same as the endpoint they are calling, e.g., `GetInventoryItem`

### Context Usage
- The first parameter of every request method should be `ctx context.Context`
- Use the context for request cancellation and timeouts

### Builder Pattern
- Request methods should utilize the Builder pattern for constructing requests
- Example:
  ```go
type TestDataBuilder struct {
data *TestData
err  []error
}

func NewTestDataBuilder() *TestDataBuilder {
return &TestDataBuilder{
data: &TestData{},
err:  make([]error, 0),
}
}

func (b *TestDataBuilder) Id(id uuid.UUID) *TestDataBuilder {
b.data.id = id
return b
}

func (b *TestDataBuilder) Page(page int) *TestDataBuilder {
if page <= 0 {
b.err = append(b.err, errors.New("page must be greater than 0"))
}
b.data.page = page
return b
}

func (b *TestDataBuilder) Build() (*TestData, error) {
if len(b.err) > 0 {
return nil, errors.Join(b.err...)
}
    err := b.requiredField()
    if len(err) > 0 {
        return nil, errors.Join(err...)
    }
    return b.data, nil
}

func (b *TestDataBuilder) requiredField() []error {
if b.data.id == uuid.Nil {
b.err = append(b.err, errors.New("id is required"))
}
return b.err
}

  ```
