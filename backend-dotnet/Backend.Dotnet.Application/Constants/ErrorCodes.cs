using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Constants
{
    public static class ErrorCodes
    {
        // Common Errors
        public const string NotFound = "NOT_FOUND";
        public const string InternalError = "INTERNAL_ERROR";
        public const string BusinessRuleViolation = "BUSINESS_RULE_VIOLATION";

        // Customer
        public const string CustomerNotFound = "CUSTOMER_NOT_FOUND";
        public const string CustomerDuplicateEmail = "CUSTOMER_DUPLICATE_EMAIL";
        public const string CustomerDuplicatePhone = "CUSTOMER_DUPLICATE_PHONE";
        public const string CustomerCreationFailed = "CUSTOMER_CREATION_FAILED";
        public const string CustomerUpdateFailed = "CUSTOMER_UPDATE_FAILED";
        public const string CustomerDeletionFailed = "CUSTOMER_DELETION_FAILED";
        public const string CustomerRestoreFailed = "CUSTOMER_RESTORE_FAILED";
        public const string CustomerAlreadyDeleted = "CUSTOMER_ALREADY_DELETED";
        public const string CustomerHasActiveVehicles = "CUSTOMER_HAS_ACTIVE_VEHICLES";

        // Vehicle
        public const string VehicleNotFound = "VEHICLE_NOT_FOUND";
        public const string VehicleDuplicateVin = "VEHICLE_DUPLICATE_VIN";
        public const string VehicleDuplicateLicensePlate = "VEHICLE_DUPLICATE_LICENSE_PLATE";
        public const string VehicleInvalidVin = "VEHICLE_INVALID_VIN";
        public const string VehicleInvalidLicensePlate = "VEHICLE_INVALID_LICENSE_PLATE";
        public const string VehicleCreationFailed = "VEHICLE_CREATION_FAILED";
        public const string VehicleUpdateFailed = "VEHICLE_UPDATE_FAILED";
        public const string VehicleDeletionFailed = "VEHICLE_DELETION_FAILED";
        public const string VehicleRestoreFailed = "VEHICLE_RESTORE_FAILED";
        public const string VehicleAlreadyDeleted = "VEHICLE_ALREADY_DELETED";
        public const string VehicleHasActiveServices = "VEHICLE_HAS_ACTIVE_SERVICES";
        public const string VehicleTransferFailed = "VEHICLE_TRANSFER_FAILED";
        public const string VehicleTransferSameCustomer = "VEHICLE_TRANSFER_SAME_CUSTOMER";

        // Vehicle Model
        public const string ModelNotFound = "MODEL_NOT_FOUND";
        public const string ModelDuplicate = "MODEL_DUPLICATE";
        public const string ModelInUse = "MODEL_IN_USE";
        public const string ModelInvalidYear = "MODEL_INVALID_YEAR";
        public const string ModelCreationFailed = "MODEL_CREATION_FAILED";
        public const string ModelUpdateFailed = "MODEL_UPDATE_FAILED";
        public const string ModelDeletionFailed = "MODEL_DELETION_FAILED";
        public const string ModelBrandNotFound = "MODEL_BRAND_NOT_FOUND";
        public const string ModelYearNotFound = "MODEL_YEAR_NOT_FOUND";

        // Reference Errors
        public const string ReferenceCustomerNotFound = "REFERENCE_CUSTOMER_NOT_FOUND";
        public const string ReferenceModelNotFound = "REFERENCE_MODEL_NOT_FOUND";
        public const string ReferenceVehicleNotFound = "REFERENCE_VEHICLE_NOT_FOUND";
        public const string ReferenceInvalid = "REFERENCE_INVALID";

        // Validation Errors
        public const string ValidationRequiredField = "VALIDATION_REQUIRED_FIELD";
        public const string ValidationMaxLength = "VALIDATION_MAX_LENGTH";
        public const string ValidationMinLength = "VALIDATION_MIN_LENGTH";
        public const string ValidationInvalidEmail = "VALIDATION_INVALID_EMAIL";
        public const string ValidationInvalidPhone = "VALIDATION_INVALID_PHONE";
        public const string ValidationInvalidDateRange = "VALIDATION_INVALID_DATE_RANGE";
        public const string ValidationFailed = "VALIDATION_FAILED";

        // Query Errors
        public const string QueryNoResults = "QUERY_NO_RESULTS";
        public const string QueryInvalidParameters = "QUERY_INVALID_PARAMETERS";
        public const string QueryFailed = "QUERY_FAILED";
    }
}
